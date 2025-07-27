package pgstore

import (
	"time"

	"github.com/google/uuid"
)

// Enums
type Role string

const (
	RolePersonal Role = "PERSONAL"
	RoleStudent  Role = "STUDENT"
)

type SchedulingStatus string

const (
	SchedulingStatusPendingConfirmation SchedulingStatus = "PENDING_CONFIRMATION"
	SchedulingStatusScheduled           SchedulingStatus = "SCHEDULED"
	SchedulingStatusInProgress          SchedulingStatus = "IN_PROGRESS"
	SchedulingStatusRescheduled         SchedulingStatus = "RESCHEDULED"
	SchedulingStatusCanceled            SchedulingStatus = "CANCELED"
	SchedulingStatusCompleted           SchedulingStatus = "COMPLETED"
	SchedulingStatusMissed              SchedulingStatus = "MISSED"
)

type SchedulingType string

const (
	SchedulingTypeOnline   SchedulingType = "ONLINE"
	SchedulingTypeInPerson SchedulingType = "IN_PERSON"
)

type Level string

const (
	LevelBeginner     Level = "BEGINNER"
	LevelIntermediary Level = "INTERMEDIARY"
	LevelAdvanced     Level = "ADVANCED"
)

type Day string

const (
	DayDom Day = "Dom"
	DaySeg Day = "Seg"
	DayTer Day = "Ter"
	DayQua Day = "Qua"
	DayQui Day = "Qui"
	DaySex Day = "Sex"
	DaySab Day = "Sab"
)

// Database table models - these are the authoritative data shapes
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	AvatarURL *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	Password  string    `json:"-" db:"password"`
	Role      Role      `json:"role" db:"role"`
	StudentID uuid.UUID `json:"student_id" db:"student_id"`
	PersonalID *uuid.UUID `json:"personal_id,omitempty" db:"personal_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Personal struct {
	ID             uuid.UUID `json:"id" db:"id"`
	PlanID         *string   `json:"plan_id,omitempty" db:"plan_id"`
	Rating         *float64  `json:"rating,omitempty" db:"rating"`
	Description    *string   `json:"description,omitempty" db:"description"`
	VideoURL       *string   `json:"video_url,omitempty" db:"video_url"`
	Experience     *string   `json:"experience,omitempty" db:"experience"`
	Specialization *string   `json:"specialization,omitempty" db:"specialization"`
	Qualifications *string   `json:"qualifications,omitempty" db:"qualifications"`
}

type Student struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	BornDate              time.Time  `json:"born_date" db:"born_date"`
	Age                   int32      `json:"age" db:"age"`
	Weight                float64    `json:"weight" db:"weight"`
	Objective             string     `json:"objective" db:"objective"`
	TrainingFrequency     string     `json:"training_frequency" db:"training_frequency"`
	DidBodybuilding       bool       `json:"did_bodybuilding" db:"did_bodybuilding"`
	MedicalCondition      *string    `json:"medical_condition,omitempty" db:"medical_condition"`
	PhysicalActivityLevel *string    `json:"physical_activity_level,omitempty" db:"physical_activity_level"`
	Observations          *string    `json:"observations,omitempty" db:"observations"`
	PersonalID            *uuid.UUID `json:"personal_id,omitempty" db:"personal_id"`
	PlanID                *string    `json:"plan_id,omitempty" db:"plan_id"`
}

type PersonalSchedule struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	WeekDay            int32     `json:"week_day" db:"week_day"`
	TimeStartInMinutes int32     `json:"time_start_in_minutes" db:"time_start_in_minutes"`
	TimeEndInMinutes   int32     `json:"time_end_in_minutes" db:"time_end_in_minutes"`
	PersonalID         uuid.UUID `json:"personal_id" db:"personal_id"`
}

type Plan struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description []string   `json:"description" db:"description"`
	Price       float64    `json:"price" db:"price"`
	PersonalID  *uuid.UUID `json:"personal_id,omitempty" db:"personal_id"`
}

type Workout struct {
	ID                       uuid.UUID  `json:"id" db:"id"`
	Name                     string     `json:"name" db:"name"`
	Description              *string    `json:"description,omitempty" db:"description"`
	Thumbnail                string     `json:"thumbnail" db:"thumbnail"`
	VideoURL                 *string    `json:"video_url,omitempty" db:"video_url"`
	RestTimeBetweenExercises *int32     `json:"rest_time_between_exercises,omitempty" db:"rest_time_between_exercises"`
	Level                    *Level     `json:"level,omitempty" db:"level"`
	WeekDays                 []Day      `json:"week_days" db:"week_days"`
	Exclusive                bool       `json:"exclusive" db:"exclusive"`
	IsTemplate               bool       `json:"is_template" db:"is_template"`
	Modality                 string     `json:"modality" db:"modality"`
	PersonalID               *uuid.UUID `json:"personal_id,omitempty" db:"personal_id"`
	StudentID                *uuid.UUID `json:"student_id,omitempty" db:"student_id"`
	PlanID                   *string    `json:"plan_id,omitempty" db:"plan_id"`
	CreatedAt                time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type ExercisesTemplate struct {
	ID                  uuid.UUID  `json:"id" db:"id"`
	Name                string     `json:"name" db:"name"`
	Thumbnail           string     `json:"thumbnail" db:"thumbnail"`
	VideoURL            string     `json:"video_url" db:"video_url"`
	Load                *int32     `json:"load,omitempty" db:"load"`
	Sets                int32      `json:"sets" db:"sets"`
	Reps                int32      `json:"reps" db:"reps"`
	RestTimeBetweenSets *int32     `json:"rest_time_between_sets,omitempty" db:"rest_time_between_sets"`
	PersonalID          *uuid.UUID `json:"personal_id,omitempty" db:"personal_id"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}

type ExercisesSetup struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	Name                string    `json:"name" db:"name"`
	Thumbnail           string    `json:"thumbnail" db:"thumbnail"`
	VideoURL            string    `json:"video_url" db:"video_url"`
	Sets                int32     `json:"sets" db:"sets"`
	Reps                int32     `json:"reps" db:"reps"`
	RestTimeBetweenSets int32     `json:"rest_time_between_sets" db:"rest_time_between_sets"`
	Load                int32     `json:"load" db:"load"`
	WorkoutID           uuid.UUID `json:"workout_id" db:"workout_id"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type Scheduling struct {
	ID          uuid.UUID        `json:"id" db:"id"`
	PersonalID  uuid.UUID        `json:"personal_id" db:"personal_id"`
	StudentID   uuid.UUID        `json:"student_id" db:"student_id"`
	WorkoutID   *uuid.UUID       `json:"workout_id,omitempty" db:"workout_id"`
	Date        time.Time        `json:"date" db:"date"`
	Type        SchedulingType   `json:"type" db:"type"`
	Status      SchedulingStatus `json:"status" db:"status"`
	StartedAt   *time.Time       `json:"started_at,omitempty" db:"started_at"`
	CompletedAt *time.Time       `json:"completed_at,omitempty" db:"completed_at"`
	StardAt     *time.Time       `json:"stard_at,omitempty" db:"stard_at"` // Note: keeping original typo for compatibility
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UserID      *uuid.UUID       `json:"user_id,omitempty" db:"user_id"`
}

type SchedulingsHistory struct {
	ID         uuid.UUID        `json:"id" db:"id"`
	ScheduleID uuid.UUID        `json:"schedule_id" db:"schedule_id"`
	UserID     uuid.UUID        `json:"user_id" db:"user_id"`
	Status     SchedulingStatus `json:"status" db:"status"`
	ChangedAt  *time.Time       `json:"changed_at,omitempty" db:"changed_at"`
	ChangedBy  string           `json:"changed_by" db:"changed_by"`
	Reason     *string          `json:"reason,omitempty" db:"reason"`
	Notes      *string          `json:"notes,omitempty" db:"notes"`
}

type Message struct {
	ID         uuid.UUID `json:"id" db:"id"`
	PersonalID uuid.UUID `json:"personal_id" db:"personal_id"`
	StudentID  uuid.UUID `json:"student_id" db:"student_id"`
	Title      string    `json:"title" db:"title"`
	Content    string    `json:"content" db:"content"`
	SentAt     time.Time `json:"sent_at" db:"sent_at"`
}

type WorkoutsHistory struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	StudentID           uuid.UUID `json:"student_id" db:"student_id"`
	WorkoutID           uuid.UUID `json:"workout_id" db:"workout_id"`
	ExecutionTime       *string   `json:"execution_time,omitempty" db:"execution_time"`
	Weight              int32     `json:"weight" db:"weight"`
	Sets                string    `json:"sets" db:"sets"`
	Reps                string    `json:"reps" db:"reps"`
	RestTime            *int32    `json:"rest_time,omitempty" db:"rest_time"`
	Thumbnail           *string   `json:"thumbnail,omitempty" db:"thumbnail"`
	TimeTotalWorkout    int32     `json:"time_total_workout" db:"time_total_workout"`
	ExerciseTitle       string    `json:"exercise_title" db:"exercise_title"`
	ExerciseID          uuid.UUID `json:"exercise_id" db:"exercise_id"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type Comment struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Content    string    `json:"content" db:"content"`
	StudentID  uuid.UUID `json:"student_id" db:"student_id"`
	PersonalID uuid.UUID `json:"personal_id" db:"personal_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type WorkoutsRating struct {
	ID         uuid.UUID `json:"id" db:"id"`
	StudentID  uuid.UUID `json:"student_id" db:"student_id"`
	WorkoutID  uuid.UUID `json:"workout_id" db:"workout_id"`
	PersonalID uuid.UUID `json:"personal_id" db:"personal_id"`
	Rating     int32     `json:"rating" db:"rating"`
	Comment    *string   `json:"comment,omitempty" db:"comment"`
	RatingDate time.Time `json:"rating_date" db:"rating_date"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type PasswordResetToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty" db:"used_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

type RefreshToken struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	Token      string     `json:"token" db:"token"`
	ExpiresAt  time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	DeviceInfo *string    `json:"device_info,omitempty" db:"device_info"`
	IPAddress  *string    `json:"ip_address,omitempty" db:"ip_address"`
}

// Additional types for messages and other features
type CreateMessageRequest struct {
	PersonalID uuid.UUID `json:"personal_id" validate:"required"`
	StudentID  uuid.UUID `json:"student_id" validate:"required"`
	Title      string    `json:"title" validate:"required,min=1,max=200"`
	Content    string    `json:"content" validate:"required,min=1,max=2000"`
}

type MessageResponse struct {
	ID         uuid.UUID `json:"id"`
	PersonalID uuid.UUID `json:"personal_id"`
	StudentID  uuid.UUID `json:"student_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sent_at"`
}
