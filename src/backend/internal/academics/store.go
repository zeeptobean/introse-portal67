package academics

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store holds the database dependency for course registration, scheduling, and academic results.
type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

// withTx runs fn in a transaction, committing when fn returns nil.
func (s *Store) withTx(ctx context.Context, fn func(pgx.Tx) error) error {
	return pgx.BeginFunc(ctx, s.db, fn)
}

const termColumns = "t.id, t.code, t.name, to_char(t.starts_on, 'YYYY-MM-DD'), to_char(t.ends_on, 'YYYY-MM-DD')"

func termDest(t *Term) []any {
	return []any{&t.ID, &t.Code, &t.Name, &t.StartsOn, &t.EndsOn}
}

func (s *Store) listTerms(ctx context.Context) ([]Term, error) {
	rows, err := s.db.Query(ctx, "SELECT "+termColumns+" FROM terms t ORDER BY t.starts_on DESC")
	if err != nil {
		return nil, fmt.Errorf("list terms: %w", err)
	}
	terms, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Term, error) {
		var t Term
		err := row.Scan(termDest(&t)...)
		return t, err
	})
	if err != nil {
		return nil, fmt.Errorf("list terms: %w", err)
	}
	return terms, nil
}

// sectionQuery selects sections with the current number of active
// enrollments. $1 is the viewing student's ID (0 for none), used to report
// that student's active enrollment in each section.
const sectionQuery = `
	SELECT s.id, ` + termColumns + `,
		c.id, c.code, c.name, c.credits, c.kind,
		s.section_code, u.id, u.full_name, s.capacity,
		(SELECT count(*) FROM enrollments e
			WHERE e.section_id = s.id AND e.status = 'enrolled')::int,
		(SELECT e.id FROM enrollments e
			WHERE e.section_id = s.id AND e.student_id = $1 AND e.status = 'enrolled')
	FROM sections s
	JOIN terms t ON t.id = s.term_id
	JOIN courses c ON c.id = s.course_id
	JOIN users u ON u.id = s.lecturer_id`

func scanSection(row pgx.CollectableRow) (Section, error) {
	var sec Section
	dest := append([]any{&sec.ID}, termDest(&sec.Term)...)
	dest = append(dest,
		&sec.Course.ID, &sec.Course.Code, &sec.Course.Name, &sec.Course.Credits, &sec.Course.Kind,
		&sec.SectionCode, &sec.Lecturer.ID, &sec.Lecturer.FullName, &sec.Capacity,
		&sec.EnrolledCount, &sec.EnrollmentID)
	if err := row.Scan(dest...); err != nil {
		return Section{}, err
	}
	sec.SeatsAvailable = max(sec.Capacity-sec.EnrolledCount, 0)
	return sec, nil
}

// listSections returns sections matching the filter with their meetings.
// lecturerID restricts the list to one lecturer's sections when non-zero.
func (s *Store) listSections(ctx context.Context, f SectionFilter, studentID, lecturerID int64) ([]Section, error) {
	var kind *string
	if f.Kind != nil {
		k := string(*f.Kind)
		kind = &k
	}
	var lecturer *int64
	if lecturerID != 0 {
		lecturer = &lecturerID
	}
	pattern := "%" + likeEscaper.Replace(strings.TrimSpace(f.Query)) + "%"
	rows, err := s.db.Query(ctx, sectionQuery+`
		WHERE ($2::bigint IS NULL OR s.term_id = $2)
			AND ($3::text IS NULL OR c.kind = $3)
			AND (c.code ILIKE $4 OR c.name ILIKE $4)
			AND ($5::bigint IS NULL OR s.lecturer_id = $5)
		ORDER BY t.starts_on DESC, c.code, s.section_code`,
		studentID, f.TermID, kind, pattern, lecturer)
	if err != nil {
		return nil, fmt.Errorf("list sections: %w", err)
	}
	sections, err := pgx.CollectRows(rows, scanSection)
	if err != nil {
		return nil, fmt.Errorf("list sections: %w", err)
	}
	if err := s.attachMeetings(ctx, sections); err != nil {
		return nil, err
	}
	return sections, nil
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func (s *Store) attachMeetings(ctx context.Context, sections []Section) error {
	if len(sections) == 0 {
		return nil
	}
	ids := make([]int64, len(sections))
	for i, sec := range sections {
		ids[i] = sec.ID
	}
	rows, err := s.db.Query(ctx, `
		SELECT section_id, id, weekday, to_char(starts_at, 'HH24:MI'), to_char(ends_at, 'HH24:MI'), room
		FROM class_meetings
		WHERE section_id = ANY($1)
		ORDER BY weekday, starts_at`, ids)
	if err != nil {
		return fmt.Errorf("list meetings: %w", err)
	}
	bySection := make(map[int64][]Meeting)
	var sectionID int64
	var m Meeting
	_, err = pgx.ForEachRow(rows, []any{&sectionID, &m.ID, &m.Weekday, &m.StartsAt, &m.EndsAt, &m.Room}, func() error {
		bySection[sectionID] = append(bySection[sectionID], m)
		return nil
	})
	if err != nil {
		return fmt.Errorf("list meetings: %w", err)
	}
	for i := range sections {
		sections[i].Meetings = bySection[sections[i].ID]
		if sections[i].Meetings == nil {
			sections[i].Meetings = []Meeting{}
		}
	}
	return nil
}

func (s *Store) sectionByID(ctx context.Context, id, studentID int64) (Section, error) {
	rows, err := s.db.Query(ctx, sectionQuery+" WHERE s.id = $2", studentID, id)
	if err != nil {
		return Section{}, fmt.Errorf("find section %d: %w", id, err)
	}
	sec, err := pgx.CollectExactlyOneRow(rows, scanSection)
	if err != nil {
		return Section{}, fmt.Errorf("find section %d: %w", id, err)
	}
	sections := []Section{sec}
	if err := s.attachMeetings(ctx, sections); err != nil {
		return Section{}, err
	}
	return sections[0], nil
}

// lockUser locks the user's row so that one student's registration changes
// run one at a time. It must be the first lock taken by those transactions.
func (s *Store) lockUser(ctx context.Context, tx pgx.Tx, id int64) (role string, active bool, err error) {
	err = tx.QueryRow(ctx, "SELECT role, is_active FROM users WHERE id = $1 FOR UPDATE", id).Scan(&role, &active)
	if err != nil {
		return "", false, fmt.Errorf("lock user %d: %w", id, err)
	}
	return role, active, nil
}

type lockedSection struct {
	CourseID   int64
	TermID     int64
	LecturerID int64
	Capacity   int32
	TermEndsOn time.Time
}

// lockSection locks the section row. Seat counting must happen while the lock
// is held so concurrent registrations cannot exceed the capacity.
func (s *Store) lockSection(ctx context.Context, tx pgx.Tx, id int64) (lockedSection, error) {
	var sec lockedSection
	err := tx.QueryRow(ctx, `
		SELECT s.course_id, s.term_id, s.lecturer_id, s.capacity, t.ends_on
		FROM sections s JOIN terms t ON t.id = s.term_id
		WHERE s.id = $1
		FOR UPDATE OF s`, id).
		Scan(&sec.CourseID, &sec.TermID, &sec.LecturerID, &sec.Capacity, &sec.TermEndsOn)
	if err != nil {
		return lockedSection{}, fmt.Errorf("lock section %d: %w", id, err)
	}
	return sec, nil
}

// enrollmentFor returns the student's enrollment row for the section, if any.
func (s *Store) enrollmentFor(ctx context.Context, tx pgx.Tx, studentID, sectionID int64) (id int64, status string, err error) {
	err = tx.QueryRow(ctx, "SELECT id, status FROM enrollments WHERE student_id = $1 AND section_id = $2",
		studentID, sectionID).Scan(&id, &status)
	if err != nil {
		return 0, "", fmt.Errorf("find enrollment: %w", err)
	}
	return id, status, nil
}

// enrolledInCourse reports whether the student actively takes another section
// of the course in the term.
func (s *Store) enrolledInCourse(ctx context.Context, tx pgx.Tx, studentID, courseID, termID, sectionID int64) (bool, error) {
	var exists bool
	err := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM enrollments e JOIN sections s ON s.id = e.section_id
			WHERE e.student_id = $1 AND e.status = 'enrolled'
				AND s.course_id = $2 AND s.term_id = $3 AND s.id <> $4
		)`, studentID, courseID, termID, sectionID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check course enrollment: %w", err)
	}
	return exists, nil
}

func (s *Store) countActive(ctx context.Context, tx pgx.Tx, sectionID int64) (int32, error) {
	var n int32
	err := tx.QueryRow(ctx, "SELECT count(*)::int FROM enrollments WHERE section_id = $1 AND status = 'enrolled'",
		sectionID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("count enrollments in section %d: %w", sectionID, err)
	}
	return n, nil
}

// scheduleConflicts returns the student's active meetings that overlap any
// meeting of the section: same weekday, overlapping times, and term date
// ranges that share at least one occurrence of that weekday. Back-to-back
// meetings do not clash.
func (s *Store) scheduleConflicts(ctx context.Context, tx pgx.Tx, studentID, sectionID int64) ([]Conflict, error) {
	rows, err := tx.Query(ctx, `
		SELECT DISTINCT c.code, s.section_code, m.weekday,
			to_char(m.starts_at, 'HH24:MI'), to_char(m.ends_at, 'HH24:MI'), m.room
		FROM class_meetings nm
		JOIN sections ns ON ns.id = nm.section_id
		JOIN terms nt ON nt.id = ns.term_id
		JOIN enrollments e ON e.student_id = $1 AND e.status = 'enrolled'
		JOIN sections s ON s.id = e.section_id
		JOIN terms t ON t.id = s.term_id
		JOIN courses c ON c.id = s.course_id
		JOIN class_meetings m ON m.section_id = s.id
		WHERE nm.section_id = $2
			AND s.id <> $2
			AND m.weekday = nm.weekday
			AND m.starts_at < nm.ends_at AND nm.starts_at < m.ends_at
			AND EXISTS (
				SELECT 1
				FROM generate_series(greatest(t.starts_on, nt.starts_on)::timestamp,
					least(t.ends_on, nt.ends_on)::timestamp, interval '1 day') AS d
				WHERE extract(isodow FROM d) = m.weekday
			)
		ORDER BY c.code, s.section_code, m.weekday`, studentID, sectionID)
	if err != nil {
		return nil, fmt.Errorf("check schedule conflicts: %w", err)
	}
	conflicts, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Conflict, error) {
		var c Conflict
		err := row.Scan(&c.CourseCode, &c.SectionCode, &c.Weekday, &c.StartsAt, &c.EndsAt, &c.Room)
		return c, err
	})
	if err != nil {
		return nil, fmt.Errorf("check schedule conflicts: %w", err)
	}
	return conflicts, nil
}

func (s *Store) insertEnrollment(ctx context.Context, tx pgx.Tx, studentID, sectionID int64) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx, "INSERT INTO enrollments (student_id, section_id) VALUES ($1, $2) RETURNING id",
		studentID, sectionID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert enrollment: %w", err)
	}
	return id, nil
}

// reactivateEnrollment turns a cancelled enrollment back into an active one
// and resets its timestamps.
func (s *Store) reactivateEnrollment(ctx context.Context, tx pgx.Tx, id int64) error {
	_, err := tx.Exec(ctx, `
		UPDATE enrollments
		SET status = 'enrolled', registered_at = now(), cancelled_at = NULL
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("reactivate enrollment %d: %w", id, err)
	}
	return nil
}

// enrollmentSection returns the section and status of the student's
// enrollment, or pgx.ErrNoRows when it does not belong to the student.
func (s *Store) enrollmentSection(ctx context.Context, tx pgx.Tx, enrollmentID, studentID int64) (sectionID int64, status string, err error) {
	err = tx.QueryRow(ctx, "SELECT section_id, status FROM enrollments WHERE id = $1 AND student_id = $2",
		enrollmentID, studentID).Scan(&sectionID, &status)
	if err != nil {
		return 0, "", fmt.Errorf("find enrollment %d: %w", enrollmentID, err)
	}
	return sectionID, status, nil
}

func (s *Store) hasGrade(ctx context.Context, tx pgx.Tx, enrollmentID int64) (bool, error) {
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM grades WHERE enrollment_id = $1)", enrollmentID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check grade for enrollment %d: %w", enrollmentID, err)
	}
	return exists, nil
}

func (s *Store) cancelEnrollment(ctx context.Context, tx pgx.Tx, id int64) error {
	_, err := tx.Exec(ctx, "UPDATE enrollments SET status = 'cancelled', cancelled_at = now() WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("cancel enrollment %d: %w", id, err)
	}
	return nil
}

const enrollmentQuery = `
	SELECT e.id, e.status, e.registered_at, e.cancelled_at, s.id, s.section_code, ` + termColumns + `,
		c.id, c.code, c.name, c.credits, c.kind, u.full_name,
		EXISTS (SELECT 1 FROM grades g WHERE g.enrollment_id = e.id)
	FROM enrollments e
	JOIN sections s ON s.id = e.section_id
	JOIN terms t ON t.id = s.term_id
	JOIN courses c ON c.id = s.course_id
	JOIN users u ON u.id = s.lecturer_id`

func scanEnrollment(row pgx.CollectableRow) (Enrollment, error) {
	var e Enrollment
	dest := []any{&e.ID, &e.Status, &e.RegisteredAt, &e.CancelledAt, &e.SectionID, &e.SectionCode}
	dest = append(dest, termDest(&e.Term)...)
	dest = append(dest, &e.Course.ID, &e.Course.Code, &e.Course.Name, &e.Course.Credits, &e.Course.Kind,
		&e.LecturerName, &e.Graded)
	err := row.Scan(dest...)
	return e, err
}

func (s *Store) enrollmentByID(ctx context.Context, id int64) (Enrollment, error) {
	rows, err := s.db.Query(ctx, enrollmentQuery+" WHERE e.id = $1", id)
	if err != nil {
		return Enrollment{}, fmt.Errorf("find enrollment %d: %w", id, err)
	}
	e, err := pgx.CollectExactlyOneRow(rows, scanEnrollment)
	if err != nil {
		return Enrollment{}, fmt.Errorf("find enrollment %d: %w", id, err)
	}
	return e, nil
}

func (s *Store) studentEnrollments(ctx context.Context, studentID int64, termID *int64) ([]Enrollment, error) {
	rows, err := s.db.Query(ctx, enrollmentQuery+`
		WHERE e.student_id = $1 AND ($2::bigint IS NULL OR s.term_id = $2)
		ORDER BY t.starts_on DESC, e.status DESC, c.code`, studentID, termID)
	if err != nil {
		return nil, fmt.Errorf("list enrollments: %w", err)
	}
	enrollments, err := pgx.CollectRows(rows, scanEnrollment)
	if err != nil {
		return nil, fmt.Errorf("list enrollments: %w", err)
	}
	return enrollments, nil
}

func (s *Store) timetable(ctx context.Context, studentID int64, termID *int64) ([]TimetableEntry, error) {
	rows, err := s.db.Query(ctx, `
		SELECT m.id, m.weekday, to_char(m.starts_at, 'HH24:MI'), to_char(m.ends_at, 'HH24:MI'), m.room,
			e.id, s.id, s.section_code, `+termColumns+`, c.code, c.name, u.full_name
		FROM enrollments e
		JOIN sections s ON s.id = e.section_id
		JOIN terms t ON t.id = s.term_id
		JOIN courses c ON c.id = s.course_id
		JOIN users u ON u.id = s.lecturer_id
		JOIN class_meetings m ON m.section_id = s.id
		WHERE e.student_id = $1 AND e.status = 'enrolled' AND ($2::bigint IS NULL OR s.term_id = $2)
		ORDER BY t.starts_on, m.weekday, m.starts_at, c.code`, studentID, termID)
	if err != nil {
		return nil, fmt.Errorf("list timetable: %w", err)
	}
	entries, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (TimetableEntry, error) {
		var te TimetableEntry
		dest := []any{&te.ID, &te.Weekday, &te.StartsAt, &te.EndsAt, &te.Room, &te.EnrollmentID, &te.SectionID, &te.SectionCode}
		dest = append(dest, termDest(&te.Term)...)
		dest = append(dest, &te.CourseCode, &te.CourseName, &te.LecturerName)
		err := row.Scan(dest...)
		return te, err
	})
	if err != nil {
		return nil, fmt.Errorf("list timetable: %w", err)
	}
	return entries, nil
}

func (s *Store) exams(ctx context.Context, studentID int64, termID *int64) ([]ExamEntry, error) {
	rows, err := s.db.Query(ctx, `
		SELECT x.id, x.label, x.starts_at, x.ends_at, x.room, s.id, s.section_code, `+termColumns+`, c.code, c.name
		FROM enrollments e
		JOIN sections s ON s.id = e.section_id
		JOIN terms t ON t.id = s.term_id
		JOIN courses c ON c.id = s.course_id
		JOIN exams x ON x.section_id = s.id
		WHERE e.student_id = $1 AND e.status = 'enrolled' AND ($2::bigint IS NULL OR s.term_id = $2)
		ORDER BY x.starts_at, c.code`, studentID, termID)
	if err != nil {
		return nil, fmt.Errorf("list exams: %w", err)
	}
	exams, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (ExamEntry, error) {
		var x ExamEntry
		dest := []any{&x.ID, &x.Label, &x.StartsAt, &x.EndsAt, &x.Room, &x.SectionID, &x.SectionCode}
		dest = append(dest, termDest(&x.Term)...)
		dest = append(dest, &x.CourseCode, &x.CourseName)
		err := row.Scan(dest...)
		return x, err
	})
	if err != nil {
		return nil, fmt.Errorf("list exams: %w", err)
	}
	return exams, nil
}

// results returns the student's active enrollments with their grade mapped
// through the term's grading scale (the band with the greatest lower bound
// not above the score).
func (s *Store) results(ctx context.Context, studentID int64) ([]CourseResult, error) {
	rows, err := s.db.Query(ctx, `
		SELECT e.id, `+termColumns+`, c.id, c.code, c.name, c.credits, c.kind, s.section_code,
			g.score::float8, b.letter_grade, b.grade_points::float8
		FROM enrollments e
		JOIN sections s ON s.id = e.section_id
		JOIN terms t ON t.id = s.term_id
		JOIN courses c ON c.id = s.course_id
		LEFT JOIN grades g ON g.enrollment_id = e.id
		LEFT JOIN LATERAL (
			SELECT gb.letter_grade, gb.grade_points
			FROM grade_bands gb
			WHERE gb.grading_scale_id = t.grading_scale_id AND gb.lower_bound <= g.score
			ORDER BY gb.lower_bound DESC
			LIMIT 1
		) b ON true
		WHERE e.student_id = $1 AND e.status = 'enrolled'
		ORDER BY t.starts_on, c.code`, studentID)
	if err != nil {
		return nil, fmt.Errorf("list results: %w", err)
	}
	results, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (CourseResult, error) {
		var r CourseResult
		dest := append([]any{&r.EnrollmentID}, termDest(&r.Term)...)
		dest = append(dest, &r.Course.ID, &r.Course.Code, &r.Course.Name, &r.Course.Credits, &r.Course.Kind,
			&r.SectionCode, &r.Score, &r.LetterGrade, &r.GradePoints)
		err := row.Scan(dest...)
		return r, err
	})
	if err != nil {
		return nil, fmt.Errorf("list results: %w", err)
	}
	return results, nil
}

// lockSectionShared returns the section's lecturer and holds a share lock so
// that enrollments cannot be cancelled while grades are written.
func (s *Store) lockSectionShared(ctx context.Context, tx pgx.Tx, sectionID int64) (lecturerID int64, err error) {
	err = tx.QueryRow(ctx, "SELECT lecturer_id FROM sections WHERE id = $1 FOR SHARE", sectionID).Scan(&lecturerID)
	if err != nil {
		return 0, fmt.Errorf("lock section %d: %w", sectionID, err)
	}
	return lecturerID, nil
}

func (s *Store) sectionLecturer(ctx context.Context, sectionID int64) (int64, error) {
	var lecturerID int64
	err := s.db.QueryRow(ctx, "SELECT lecturer_id FROM sections WHERE id = $1", sectionID).Scan(&lecturerID)
	if err != nil {
		return 0, fmt.Errorf("find section %d: %w", sectionID, err)
	}
	return lecturerID, nil
}

func (s *Store) roster(ctx context.Context, sectionID int64) ([]RosterEntry, error) {
	rows, err := s.db.Query(ctx, `
		SELECT e.id, u.id, u.institutional_id, u.full_name, u.email, g.score::float8, g.updated_at
		FROM enrollments e
		JOIN users u ON u.id = e.student_id
		LEFT JOIN grades g ON g.enrollment_id = e.id
		WHERE e.section_id = $1 AND e.status = 'enrolled'
		ORDER BY u.full_name, u.id`, sectionID)
	if err != nil {
		return nil, fmt.Errorf("list roster: %w", err)
	}
	roster, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (RosterEntry, error) {
		var r RosterEntry
		err := row.Scan(&r.EnrollmentID, &r.StudentID, &r.InstitutionalID, &r.FullName, &r.Email, &r.Score, &r.GradedAt)
		return r, err
	})
	if err != nil {
		return nil, fmt.Errorf("list roster: %w", err)
	}
	return roster, nil
}

// activeEnrollmentIDs returns which of the given enrollments are active in the section.
func (s *Store) activeEnrollmentIDs(ctx context.Context, tx pgx.Tx, sectionID int64, ids []int64) (map[int64]bool, error) {
	rows, err := tx.Query(ctx, `
		SELECT id FROM enrollments
		WHERE section_id = $1 AND status = 'enrolled' AND id = ANY($2)`, sectionID, ids)
	if err != nil {
		return nil, fmt.Errorf("find section enrollments: %w", err)
	}
	found, err := pgx.CollectRows(rows, pgx.RowTo[int64])
	if err != nil {
		return nil, fmt.Errorf("find section enrollments: %w", err)
	}
	active := make(map[int64]bool, len(found))
	for _, id := range found {
		active[id] = true
	}
	return active, nil
}

func (s *Store) upsertGrade(ctx context.Context, tx pgx.Tx, enrollmentID, lecturerID int64, score float64) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO grades (enrollment_id, score, graded_by, updated_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (enrollment_id) DO UPDATE
		SET score = EXCLUDED.score, graded_by = EXCLUDED.graded_by, updated_at = EXCLUDED.updated_at`,
		enrollmentID, score, lecturerID)
	if err != nil {
		return fmt.Errorf("save grade for enrollment %d: %w", enrollmentID, err)
	}
	return nil
}
