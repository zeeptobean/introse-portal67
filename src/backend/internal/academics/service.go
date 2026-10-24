package academics

import (
	"context"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

// campus is the time zone of class meetings and term dates.
var campus = loadCampus()

// loadCampus falls back to a fixed UTC+7 zone when tzdata is unavailable.
func loadCampus() *time.Location {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return time.FixedZone("ICT", 7*60*60)
	}
	return loc
}

// Service owns business rules and transaction boundaries for course registration, scheduling, and academic results.
type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) Terms(ctx context.Context) ([]Term, error) {
	return s.store.listTerms(ctx)
}

// Sections lists offered sections. studentID marks the student's own active
// enrollments and may be zero.
func (s *Service) Sections(ctx context.Context, f SectionFilter, studentID int64) ([]Section, error) {
	return s.store.listSections(ctx, f, studentID, 0)
}

func (s *Service) LecturerSections(ctx context.Context, lecturerID int64) ([]Section, error) {
	return s.store.listSections(ctx, SectionFilter{}, 0, lecturerID)
}

// Register enrolls the student in the section. Every registration and
// cancellation locks the student row first and then the section row; see
// docs/design/data-model.md.
func (s *Service) Register(ctx context.Context, studentID, sectionID int64) (Enrollment, error) {
	var enrollmentID int64
	err := s.store.withTx(ctx, func(tx pgx.Tx) error {
		role, active, err := s.store.lockUser(ctx, tx, studentID)
		if err != nil {
			return err
		}
		if role != "student" || !active {
			return ErrNotStudent
		}
		sec, err := s.store.lockSection(ctx, tx, sectionID)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrSectionNotFound
		}
		if err != nil {
			return err
		}
		if sec.TermEndsOn.Before(today()) {
			return ErrTermEnded
		}

		existingID, status, err := s.store.enrollmentFor(ctx, tx, studentID, sectionID)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
		case err != nil:
			return err
		case status == StatusEnrolled:
			return ErrAlreadyEnrolled
		}

		taken, err := s.store.enrolledInCourse(ctx, tx, studentID, sec.CourseID, sec.TermID, sectionID)
		if err != nil {
			return err
		}
		if taken {
			return ErrCourseTaken
		}
		enrolled, err := s.store.countActive(ctx, tx, sectionID)
		if err != nil {
			return err
		}
		if enrolled >= sec.Capacity {
			return ErrSectionFull
		}
		conflicts, err := s.store.scheduleConflicts(ctx, tx, studentID, sectionID)
		if err != nil {
			return err
		}
		if len(conflicts) > 0 {
			return &ScheduleConflictError{Conflicts: conflicts}
		}

		if existingID != 0 {
			enrollmentID = existingID
			return s.store.reactivateEnrollment(ctx, tx, existingID)
		}
		enrollmentID, err = s.store.insertEnrollment(ctx, tx, studentID, sectionID)
		return err
	})
	if err != nil {
		return Enrollment{}, err
	}
	return s.store.enrollmentByID(ctx, enrollmentID)
}

// Cancel cancels the student's active enrollment. Graded enrollments cannot
// be cancelled.
func (s *Service) Cancel(ctx context.Context, studentID, enrollmentID int64) (Enrollment, error) {
	err := s.store.withTx(ctx, func(tx pgx.Tx) error {
		if _, _, err := s.store.lockUser(ctx, tx, studentID); err != nil {
			return err
		}
		sectionID, _, err := s.store.enrollmentSection(ctx, tx, enrollmentID, studentID)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrEnrollmentNotFound
		}
		if err != nil {
			return err
		}
		sec, err := s.store.lockSection(ctx, tx, sectionID)
		if err != nil {
			return err
		}
		if sec.TermEndsOn.Before(today()) {
			return ErrTermEnded
		}
		// Re-read the status now that the section lock is held.
		_, status, err := s.store.enrollmentSection(ctx, tx, enrollmentID, studentID)
		if err != nil {
			return err
		}
		if status != StatusEnrolled {
			return ErrAlreadyCancelled
		}
		graded, err := s.store.hasGrade(ctx, tx, enrollmentID)
		if err != nil {
			return err
		}
		if graded {
			return ErrAlreadyGraded
		}
		return s.store.cancelEnrollment(ctx, tx, enrollmentID)
	})
	if err != nil {
		return Enrollment{}, err
	}
	return s.store.enrollmentByID(ctx, enrollmentID)
}

func (s *Service) Enrollments(ctx context.Context, studentID int64, termID *int64) ([]Enrollment, error) {
	return s.store.studentEnrollments(ctx, studentID, termID)
}

// Timetable returns the meetings of the student's active enrollments and
// marks meetings that overlap each other.
func (s *Service) Timetable(ctx context.Context, studentID int64, termID *int64) ([]TimetableEntry, error) {
	entries, err := s.store.timetable(ctx, studentID, termID)
	if err != nil {
		return nil, err
	}
	for i := range entries {
		entries[i].ConflictsWith = []string{}
	}
	for i := range entries {
		for j := i + 1; j < len(entries); j++ {
			a, b := &entries[i], &entries[j]
			if a.SectionID == b.SectionID || !meetingsOverlap(a, b) {
				continue
			}
			a.ConflictsWith = append(a.ConflictsWith, strconv.FormatInt(b.ID, 10))
			b.ConflictsWith = append(b.ConflictsWith, strconv.FormatInt(a.ID, 10))
		}
	}
	return entries, nil
}

// meetingsOverlap applies the weekly clash rule: same weekday, overlapping
// times, and term dates that share at least one occurrence of that weekday.
func meetingsOverlap(a, b *TimetableEntry) bool {
	if a.Weekday != b.Weekday || a.StartsAt >= b.EndsAt || b.StartsAt >= a.EndsAt {
		return false
	}
	from := max(a.Term.StartsOn, b.Term.StartsOn)
	to := min(a.Term.EndsOn, b.Term.EndsOn)
	start, err1 := time.Parse(time.DateOnly, from)
	end, err2 := time.Parse(time.DateOnly, to)
	if err1 != nil || err2 != nil || end.Before(start) {
		return false
	}
	for d := start; !d.After(end) && d.Sub(start) < 7*24*time.Hour; d = d.AddDate(0, 0, 1) {
		if isoWeekday(d) == a.Weekday {
			return true
		}
	}
	return false
}

func isoWeekday(t time.Time) int16 {
	if t.Weekday() == time.Sunday {
		return 7
	}
	return int16(t.Weekday())
}

// Exams returns the exams of the student's active enrollments and marks
// exams whose times overlap.
func (s *Service) Exams(ctx context.Context, studentID int64, termID *int64) ([]ExamEntry, error) {
	exams, err := s.store.exams(ctx, studentID, termID)
	if err != nil {
		return nil, err
	}
	for i := range exams {
		exams[i].ConflictsWith = []string{}
	}
	for i := range exams {
		for j := i + 1; j < len(exams); j++ {
			a, b := &exams[i], &exams[j]
			if a.StartsAt.Before(b.EndsAt) && b.StartsAt.Before(a.EndsAt) {
				a.ConflictsWith = append(a.ConflictsWith, strconv.FormatInt(b.ID, 10))
				b.ConflictsWith = append(b.ConflictsWith, strconv.FormatInt(a.ID, 10))
			}
		}
	}
	return exams, nil
}

// Results groups the student's courses by term and computes semester and
// cumulative GPA as sum(points * credits) / sum(graded credits). Ungraded
// courses are excluded; a GPA without graded credits is null.
func (s *Service) Results(ctx context.Context, studentID int64) (Results, error) {
	rows, err := s.store.results(ctx, studentID)
	if err != nil {
		return Results{}, err
	}
	out := Results{Terms: []TermResults{}}
	var termPoints []float64
	var totalPoints float64
	for _, r := range rows {
		if len(out.Terms) == 0 || out.Terms[len(out.Terms)-1].Term.ID != r.Term.ID {
			out.Terms = append(out.Terms, TermResults{Term: r.Term, Courses: []CourseResult{}})
			termPoints = append(termPoints, 0)
		}
		tr := &out.Terms[len(out.Terms)-1]
		tr.Courses = append(tr.Courses, r)
		if r.GradePoints != nil {
			tr.GradedCredits += int(r.Course.Credits)
			out.GradedCredits += int(r.Course.Credits)
			points := *r.GradePoints * float64(r.Course.Credits)
			termPoints[len(termPoints)-1] += points
			totalPoints += points
		}
	}
	for i := range out.Terms {
		out.Terms[i].GPA = gpa(termPoints[i], out.Terms[i].GradedCredits)
	}
	out.CumulativeGPA = gpa(totalPoints, out.GradedCredits)
	return out, nil
}

func gpa(points float64, credits int) *float64 {
	if credits == 0 {
		return nil
	}
	v := math.Round(points/float64(credits)*100) / 100
	return &v
}

// Roster lists the students actively enrolled in the lecturer's section.
func (s *Service) Roster(ctx context.Context, lecturerID, sectionID int64) ([]RosterEntry, error) {
	assigned, err := s.store.sectionLecturer(ctx, sectionID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrSectionNotFound
	}
	if err != nil {
		return nil, err
	}
	if assigned != lecturerID {
		return nil, ErrNotAssigned
	}
	return s.store.roster(ctx, sectionID)
}

// SaveGrades records scores for active enrollments of the lecturer's section.
// Either all grades are saved or none are.
func (s *Service) SaveGrades(ctx context.Context, lecturerID, sectionID int64, grades []GradeInput) ([]RosterEntry, error) {
	if len(grades) == 0 {
		return nil, &GradeInputError{Message: "at least one grade is required"}
	}
	ids := make([]int64, len(grades))
	for i, g := range grades {
		if g.EnrollmentID <= 0 {
			return nil, &GradeInputError{Message: fmt.Sprintf("grade %d: enrollment_id is required", i+1)}
		}
		if slices.Contains(ids[:i], g.EnrollmentID) {
			return nil, &GradeInputError{Message: fmt.Sprintf("grade %d: duplicate enrollment %d", i+1, g.EnrollmentID)}
		}
		if err := validateScore(g.Score); err != nil {
			return nil, &GradeInputError{Message: fmt.Sprintf("grade %d: %v", i+1, err)}
		}
		ids[i] = g.EnrollmentID
	}

	err := s.store.withTx(ctx, func(tx pgx.Tx) error {
		assigned, err := s.store.lockSectionShared(ctx, tx, sectionID)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrSectionNotFound
		}
		if err != nil {
			return err
		}
		if assigned != lecturerID {
			return ErrNotAssigned
		}
		active, err := s.store.activeEnrollmentIDs(ctx, tx, sectionID, ids)
		if err != nil {
			return err
		}
		for _, g := range grades {
			if !active[g.EnrollmentID] {
				return &GradeInputError{Message: fmt.Sprintf("enrollment %d is not an active enrollment of this section", g.EnrollmentID)}
			}
			if err := s.store.upsertGrade(ctx, tx, g.EnrollmentID, lecturerID, g.Score); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.store.roster(ctx, sectionID)
}

func validateScore(score float64) error {
	if math.IsNaN(score) || score < 0 || score > 10 {
		return errors.New("score must be between 0 and 10")
	}
	if math.Abs(score*100-math.Round(score*100)) > 1e-6 {
		return errors.New("score must have at most two decimal places")
	}
	return nil
}

// today returns the current campus date at midnight UTC, matching how pgx
// decodes DATE values.
func today() time.Time {
	y, m, d := time.Now().In(campus).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
