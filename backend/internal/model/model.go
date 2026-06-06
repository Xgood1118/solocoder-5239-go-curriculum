package model

import "time"

type DifficultyLevel int

const (
	Level1 DifficultyLevel = 1
	Level2 DifficultyLevel = 2
	Level3 DifficultyLevel = 3
	Level4 DifficultyLevel = 4
	Level5 DifficultyLevel = 5
)

type PlanStatus string

const (
	StatusDraft     PlanStatus = "draft"
	StatusPublished PlanStatus = "published"
	StatusExecuting PlanStatus = "executing"
	StatusCompleted PlanStatus = "completed"
	StatusArchived  PlanStatus = "archived"
)

type RevisionStatus string

const (
	RevisionPending  RevisionStatus = "pending"
	RevisionApproved RevisionStatus = "approved"
	RevisionRejected RevisionStatus = "rejected"
)

type RecordType string

const (
	RecordNormal   RecordType = "normal"
	RecordSub      RecordType = "substitute"
	RecordStudy    RecordType = "study_self"
	RecordPE       RecordType = "pe"
	RecordMeeting  RecordType = "meeting"
	RecordOccupied RecordType = "occupied"
)

type KnowledgePoint struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Difficulty  DifficultyLevel `json:"difficulty"`
	Prerequisites []string      `json:"prerequisites"`
	IsKeyPoint  bool            `json:"is_key_point"`
	IsDifficult bool            `json:"is_difficult"`
	IsExamPoint bool            `json:"is_exam_point"`
	EstimatedLessons int        `json:"estimated_lessons"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type Lesson struct {
	ID               string        `json:"id"`
	UnitID           string        `json:"unit_id"`
	Title            string        `json:"title"`
	Content          string        `json:"content"`
	KnowledgePointIDs []string     `json:"knowledge_point_ids"`
	PlanWeek         int           `json:"plan_week"`
	PlanDayIndex     int           `json:"plan_day_index"`
	OrderInUnit      int           `json:"order_in_unit"`
	DurationMinutes  int           `json:"duration_minutes"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type Unit struct {
	ID               string        `json:"id"`
	SemesterID       string        `json:"semester_id"`
	Title            string        `json:"title"`
	Description      string        `json:"description"`
	KnowledgePointIDs []string     `json:"knowledge_point_ids"`
	LessonCount      int           `json:"lesson_count"`
	Lessons          []Lesson      `json:"lessons,omitempty"`
	SuggestedMinLessons int        `json:"suggested_min_lessons"`
	SuggestedMaxLessons int        `json:"suggested_max_lessons"`
	OrderIndex       int           `json:"order_index"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type Semester struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Grade        string     `json:"grade"`
	Subject      string     `json:"subject"`
	TotalWeeks   int        `json:"total_weeks"`
	LessonsPerDay int       `json:"lessons_per_day"`
	LessonMinutes int        `json:"lesson_minutes"`
	Units        []Unit     `json:"units,omitempty"`
	Status       PlanStatus `json:"status"`
	Version      string     `json:"version"`
	ParentPlanID string     `json:"parent_plan_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CreatedBy    string     `json:"created_by"`
}

type Teacher struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Subject  string   `json:"subject"`
	ClassIDs []string `json:"class_ids"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Class struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Grade      string   `json:"grade"`
	StudentCount int     `json:"student_count"`
	TeacherIDs []string `json:"teacher_ids"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ClassPlan struct {
	ID            string     `json:"id"`
	ClassID       string     `json:"class_id"`
	Subject       string     `json:"subject"`
	SemesterID    string     `json:"semester_id"`
	CurrentTeacherID string   `json:"current_teacher_id"`
	Status        PlanStatus `json:"status"`
	IsLocked      bool       `json:"is_locked"`
	StartDate     time.Time  `json:"start_date"`
	EndDate       time.Time  `json:"end_date"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ProgressRecord struct {
	ID              string     `json:"id"`
	ClassPlanID     string     `json:"class_plan_id"`
	LessonID        string     `json:"lesson_id"`
	RecordType      RecordType `json:"record_type"`
	ActualWeek      int        `json:"actual_week"`
	ActualDate      time.Time  `json:"actual_date"`
	TeacherID       string     `json:"teacher_id"`
	IsSubstitute    bool       `json:"is_substitute"`
	SubstituteNote  string     `json:"substitute_note,omitempty"`
	Notes           string     `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type ProgressSummary struct {
	ClassPlanID        string  `json:"class_plan_id"`
	TotalLessons       int     `json:"total_lessons"`
	CompletedLessons   int     `json:"completed_lessons"`
	PlanLessonsByWeek  int     `json:"plan_lessons_by_week"`
	Deviation          int     `json:"deviation"`
	DeviationLevel     string  `json:"deviation_level"`
	CurrentWeek        int     `json:"current_week"`
	KeyCoverageRate    float64 `json:"key_coverage_rate"`
	DifficultCoverageRate float64 `json:"difficult_coverage_rate"`
	ExamCoverageRate   float64 `json:"exam_coverage_rate"`
}

type RevisionRequest struct {
	ID            string         `json:"id"`
	OriginalPlanID string        `json:"original_plan_id"`
	NewPlanID     string         `json:"new_plan_id,omitempty"`
	Title         string         `json:"title"`
	Reason        string         `json:"reason"`
	ApplicantID   string         `json:"applicant_id"`
	Status        RevisionStatus `json:"status"`
	ApproverID    string         `json:"approver_id,omitempty"`
	ApprovalNote  string         `json:"approval_note,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	ApprovedAt    *time.Time     `json:"approved_at,omitempty"`
}

type CycleDetectionResult struct {
	HasCycle bool     `json:"has_cycle"`
	Cycle    []string `json:"cycle,omitempty"`
}

type TopoSortResult struct {
	Sorted []string              `json:"sorted"`
	Levels map[string]int        `json:"levels"`
	Cycle  *CycleDetectionResult `json:"cycle,omitempty"`
}

type LessonSuggestion struct {
	MinLessons int     `json:"min_lessons"`
	MaxLessons int     `json:"max_lessons"`
	Reason     string  `json:"reason"`
}

type GanttBar struct {
	UnitID    string `json:"unit_id"`
	UnitTitle string `json:"unit_title"`
	StartWeek int    `json:"start_week"`
	EndWeek   int    `json:"end_week"`
}

type GanttCompareData struct {
	ClassID   string     `json:"class_id"`
	ClassName string     `json:"class_name"`
	Bars      []GanttBar `json:"bars"`
}

type CoverageItem struct {
	UnitID   string  `json:"unit_id"`
	UnitName string  `json:"unit_name"`
	Rate     float64 `json:"rate"`
	IsPassed bool    `json:"is_passed"`
}

type CoverageReport struct {
	KeyPoints       []CoverageItem `json:"key_points"`
	DifficultPoints []CoverageItem `json:"difficult_points"`
	ExamPoints      []CoverageItem `json:"exam_points"`
	UnpassedUnits   []string       `json:"unpassed_units"`
}

type ExcelReportRow struct {
	ClassName     string  `json:"class_name"`
	Subject       string  `json:"subject"`
	PlanProgress  float64 `json:"plan_progress"`
	ActualProgress float64 `json:"actual_progress"`
	Deviation     int     `json:"deviation"`
	IsUpToStandard bool   `json:"is_up_to_standard"`
}
