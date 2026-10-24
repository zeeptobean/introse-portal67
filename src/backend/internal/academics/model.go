package academics

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type CourseKind string

const (
	KindCourse           CourseKind = "course"
	KindSpecializedTopic CourseKind = "specialized_topic"
)

const (
	StatusEnrolled  = "enrolled"
	StatusCancelled = "cancelled"
)

var (
	ErrSectionNotFound    = errors.New("section not found")
	ErrEnrollmentNotFound = errors.New("enrollment not found")
	ErrNotStudent         = errors.New("only active students can register")
	ErrAlreadyEnrolled    = errors.New("already registered for this section")
	ErrCourseTaken        = errors.New("already registered for another section of this course in the same term")
	ErrSectionFull        = errors.New("section is full")
	ErrTermEnded          = errors.New("the term of this section has ended")
	ErrAlreadyCancelled   = errors.New("enrollment is already cancelled")
	ErrAlreadyGraded      = errors.New("cannot cancel an enrollment that has a grade")
	ErrNotAssigned        = errors.New("you are not the lecturer of this section")
)

// ScheduleConflictError lists the existing meetings that clash with the
// section a student tried to register for.
type ScheduleConflictError struct {
	Conflicts []Conflict
}

func (e *ScheduleConflictError) Error() string {
	parts := make([]string, len(e.Conflicts))
	for i, c := range e.Conflicts {
		parts[i] = fmt.Sprintf("%s section %s (%s %s-%s)",
			c.CourseCode, c.SectionCode, weekdayNames[c.Weekday], c.StartsAt, c.EndsAt)
	}
	return "timetable conflict with " + strings.Join(parts, ", ")
}

// GradeInputError reports a grade entry that cannot be saved.
type GradeInputError struct {
	Message string
}

func (e *GradeInputError) Error() string { return e.Message }

var weekdayNames = map[int16]string{1: "Mon", 2: "Tue", 3: "Wed", 4: "Thu", 5: "Fri", 6: "Sat", 7: "Sun"}

type Term struct {
	ID       int64  `json:"id,string"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	StartsOn string `json:"starts_on"`
	EndsOn   string `json:"ends_on"`
}

type Course struct {
	ID      int64      `json:"id,string"`
	Code    string     `json:"code"`
	Name    string     `json:"name"`
	Credits int16      `json:"credits"`
	Kind    CourseKind `json:"kind"`
}

type Lecturer struct {
	ID       int64  `json:"id,string"`
	FullName string `json:"full_name"`
}

// Meeting is one weekly class slot. Times are "HH:MM" in campus local time.
type Meeting struct {
	ID       int64   `json:"id,string"`
	Weekday  int16   `json:"weekday"`
	StartsAt string  `json:"starts_at"`
	EndsAt   string  `json:"ends_at"`
	Room     *string `json:"room"`
}

type Section struct {
	ID             int64     `json:"id,string"`
	Term           Term      `json:"term"`
	Course         Course    `json:"course"`
	SectionCode    string    `json:"section_code"`
	Lecturer       Lecturer  `json:"lecturer"`
	Capacity       int32     `json:"capacity"`
	EnrolledCount  int32     `json:"enrolled_count"`
	SeatsAvailable int32     `json:"seats_available"`
	Meetings       []Meeting `json:"meetings"`
	// EnrollmentID is set when the current student has an active enrollment.
	EnrollmentID *int64 `json:"enrollment_id,string,omitempty"`
}

type SectionFilter struct {
	TermID *int64
	Kind   *CourseKind
	Query  string
}

type Enrollment struct {
	ID           int64      `json:"id,string"`
	Status       string     `json:"status"`
	RegisteredAt time.Time  `json:"registered_at"`
	CancelledAt  *time.Time `json:"cancelled_at"`
	SectionID    int64      `json:"section_id,string"`
	SectionCode  string     `json:"section_code"`
	Term         Term       `json:"term"`
	Course       Course     `json:"course"`
	LecturerName string     `json:"lecturer_name"`
	Graded       bool       `json:"graded"`
}

// TimetableEntry is a meeting of one of the student's active enrollments.
// ConflictsWith holds the IDs of other entries that overlap it.
type TimetableEntry struct {
	Meeting
	EnrollmentID  int64    `json:"enrollment_id,string"`
	SectionID     int64    `json:"section_id,string"`
	SectionCode   string   `json:"section_code"`
	Term          Term     `json:"term"`
	CourseCode    string   `json:"course_code"`
	CourseName    string   `json:"course_name"`
	LecturerName  string   `json:"lecturer_name"`
	ConflictsWith []string `json:"conflicts_with"`
}

type ExamEntry struct {
	ID            int64     `json:"id,string"`
	Label         string    `json:"label"`
	StartsAt      time.Time `json:"starts_at"`
	EndsAt        time.Time `json:"ends_at"`
	Room          *string   `json:"room"`
	SectionID     int64     `json:"section_id,string"`
	SectionCode   string    `json:"section_code"`
	Term          Term      `json:"term"`
	CourseCode    string    `json:"course_code"`
	CourseName    string    `json:"course_name"`
	ConflictsWith []string  `json:"conflicts_with"`
}

// Conflict describes an existing meeting that blocks a registration.
type Conflict struct {
	CourseCode  string  `json:"course_code"`
	SectionCode string  `json:"section_code"`
	Weekday     int16   `json:"weekday"`
	StartsAt    string  `json:"starts_at"`
	EndsAt      string  `json:"ends_at"`
	Room        *string `json:"room"`
}

// CourseResult is one active enrollment with its grade, if any.
type CourseResult struct {
	EnrollmentID int64    `json:"enrollment_id,string"`
	Term         Term     `json:"-"`
	Course       Course   `json:"course"`
	SectionCode  string   `json:"section_code"`
	Score        *float64 `json:"score"`
	LetterGrade  *string  `json:"letter_grade"`
	GradePoints  *float64 `json:"grade_points"`
}

type TermResults struct {
	Term          Term           `json:"term"`
	Courses       []CourseResult `json:"courses"`
	GradedCredits int            `json:"graded_credits"`
	GPA           *float64       `json:"gpa"`
}

type Results struct {
	Terms         []TermResults `json:"terms"`
	GradedCredits int           `json:"graded_credits"`
	CumulativeGPA *float64      `json:"cumulative_gpa"`
}

// RosterEntry is a student actively enrolled in a lecturer's section.
type RosterEntry struct {
	EnrollmentID    int64      `json:"enrollment_id,string"`
	StudentID       int64      `json:"student_id,string"`
	InstitutionalID *string    `json:"institutional_id"`
	FullName        string     `json:"full_name"`
	Email           string     `json:"email"`
	Score           *float64   `json:"score"`
	GradedAt        *time.Time `json:"graded_at"`
}

type GradeInput struct {
	EnrollmentID int64   `json:"enrollment_id,string"`
	Score        float64 `json:"score"`
}
