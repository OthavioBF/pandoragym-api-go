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
	RoleAdmin    Role = "ADMIN"
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
	ID         uuid.UUID  `json:"id" db:"id"`
	Name       string     `json:"name" db:"name"`
	Email      string     `json:"email" db:"email"`
	Phone      string     `json:"phone" db:"phone"`
	AvatarURL  *string    `json:"avatarUrl,omitempty" db:"avatar_url"`
	Password   string     `json:"-" db:"password"`
	Role       Role       `json:"role" db:"role"`
	StudentID  uuid.UUID  `json:"studentId" db:"student_id"`
	PersonalID *uuid.UUID `json:"personalId,omitempty" db:"personal_id"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt" db:"updated_at"`
}

type Personal struct {
	ID             uuid.UUID `json:"id" db:"id"`
	PlanID         *string   `json:"planId,omitempty" db:"plan_id"`
	Rating         *float64  `json:"rating,omitempty" db:"rating"`
	Description    *string   `json:"description,omitempty" db:"description"`
	VideoURL       *string   `json:"videoUrl,omitempty" db:"video_url"`
	Experience     *string   `json:"experience,omitempty" db:"experience"`
	Specialization *string   `json:"specialization,omitempty" db:"specialization"`
	Qualifications *string   `json:"qualifications,omitempty" db:"qualifications"`
}

type Student struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	BornDate              time.Time  `json:"bornDate" db:"born_date"`
	Age                   int32      `json:"age" db:"age"`
	Weight                float64    `json:"weight" db:"weight"`
	Objective             string     `json:"objective" db:"objective"`
	TrainingFrequency     string     `json:"trainingFrequency" db:"training_frequency"`
	DidBodybuilding       bool       `json:"didBodybuilding" db:"did_bodybuilding"`
	MedicalCondition      *string    `json:"medicalCondition,omitempty" db:"medical_condition"`
	PhysicalActivityLevel *string    `json:"physicalActivityLevel,omitempty" db:"physical_activity_level"`
	Observations          *string    `json:"observations,omitempty" db:"observations"`
	PersonalID            *uuid.UUID `json:"personalId,omitempty" db:"personal_id"`
	PlanID                *string    `json:"planId,omitempty" db:"plan_id"`
}

type PersonalSchedule struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	WeekDay            int32     `json:"weekDay" db:"week_day"`
	TimeStartInMinutes int32     `json:"timeStartInMinutes" db:"time_start_in_minutes"`
	TimeEndInMinutes   int32     `json:"timeEndInMinutes" db:"time_end_in_minutes"`
	PersonalID         uuid.UUID `json:"personalId" db:"personal_id"`
}

type Plan struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description []string   `json:"description" db:"description"`
	Price       float64    `json:"price" db:"price"`
	PersonalID  *uuid.UUID `json:"personalId,omitempty" db:"personal_id"`
}

type Workout struct {
	ID                       uuid.UUID  `json:"id" db:"id"`
	Name                     string     `json:"name" db:"name"`
	Description              *string    `json:"description,omitempty" db:"description"`
	Thumbnail                string     `json:"thumbnail" db:"thumbnail"`
	VideoURL                 *string    `json:"videoUrl,omitempty" db:"video_url"`
	RestTimeBetweenExercises *int32     `json:"restTimeBetweenExercises,omitempty" db:"rest_time_between_exercises"`
	Level                    *Level     `json:"level,omitempty" db:"level"`
	WeekDays                 []Day      `json:"weekDays" db:"week_days"`
	Exclusive                bool       `json:"exclusive" db:"exclusive"`
	IsTemplate               bool       `json:"isTemplate" db:"is_template"`
	Modality                 string     `json:"modality" db:"modality"`
	PersonalID               *uuid.UUID `json:"personalId,omitempty" db:"personal_id"`
	StudentID                *uuid.UUID `json:"studentId,omitempty" db:"student_id"`
	PlanID                   *string    `json:"planId,omitempty" db:"plan_id"`
	CreatedAt                time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt                time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt                *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}

type ExercisesTemplate struct {
	ID                  uuid.UUID  `json:"id" db:"id"`
	Name                string     `json:"name" db:"name"`
	Thumbnail           string     `json:"thumbnail" db:"thumbnail"`
	VideoURL            string     `json:"videoUrl" db:"video_url"`
	Load                *int32     `json:"load,omitempty" db:"load"`
	Sets                int32      `json:"sets" db:"sets"`
	Reps                int32      `json:"reps" db:"reps"`
	RestTimeBetweenSets *int32     `json:"restTimeBetweenSets,omitempty" db:"rest_time_between_sets"`
	PersonalID          *uuid.UUID `json:"personalId,omitempty" db:"personal_id"`
	CreatedAt           time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time  `json:"updatedAt" db:"updated_at"`
}

type ExercisesSetup struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	Name                string    `json:"name" db:"name"`
	Thumbnail           string    `json:"thumbnail" db:"thumbnail"`
	VideoURL            string    `json:"videoUrl" db:"video_url"`
	Sets                int32     `json:"sets" db:"sets"`
	Reps                int32     `json:"reps" db:"reps"`
	RestTimeBetweenSets int32     `json:"restTimeBetweenSets" db:"rest_time_between_sets"`
	Load                int32     `json:"load" db:"load"`
	WorkoutID           uuid.UUID `json:"workoutId" db:"workout_id"`
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time `json:"updatedAt" db:"updated_at"`
}

type Scheduling struct {
	ID          uuid.UUID        `json:"id" db:"id"`
	PersonalID  uuid.UUID        `json:"personalId" db:"personal_id"`
	StudentID   uuid.UUID        `json:"studentId" db:"student_id"`
	WorkoutID   *uuid.UUID       `json:"workoutId,omitempty" db:"workout_id"`
	Date        time.Time        `json:"date" db:"date"`
	Type        SchedulingType   `json:"type" db:"type"`
	Status      SchedulingStatus `json:"status" db:"status"`
	StartedAt   *time.Time       `json:"startedAt,omitempty" db:"started_at"`
	CompletedAt *time.Time       `json:"completedAt,omitempty" db:"completed_at"`
	StardAt     *time.Time       `json:"stardAt,omitempty" db:"stard_at"` // Note: keeping original typo for compatibility
	CreatedAt   time.Time        `json:"createdAt" db:"created_at"`
	UserID      *uuid.UUID       `json:"userId,omitempty" db:"user_id"`
}

type SchedulingsHistory struct {
	ID         uuid.UUID        `json:"id" db:"id"`
	ScheduleID uuid.UUID        `json:"scheduleId" db:"schedule_id"`
	UserID     uuid.UUID        `json:"userId" db:"user_id"`
	Status     SchedulingStatus `json:"status" db:"status"`
	ChangedAt  *time.Time       `json:"changedAt,omitempty" db:"changed_at"`
	ChangedBy  string           `json:"changedBy" db:"changed_by"`
	Reason     *string          `json:"reason,omitempty" db:"reason"`
	Notes      *string          `json:"notes,omitempty" db:"notes"`
}

type Message struct {
	ID         uuid.UUID `json:"id" db:"id"`
	PersonalID uuid.UUID `json:"personalId" db:"personal_id"`
	StudentID  uuid.UUID `json:"studentId" db:"student_id"`
	Title      string    `json:"title" db:"title"`
	Content    string    `json:"content" db:"content"`
	SentAt     time.Time `json:"sentAt" db:"sent_at"`
}

type WorkoutsHistory struct {
	ID                uuid.UUID `json:"id" db:"id"`
	StudentID         uuid.UUID `json:"studentId" db:"student_id"`
	WorkoutID         uuid.UUID `json:"workoutId" db:"workout_id"`
	ExecutionTime     *string   `json:"executionTime,omitempty" db:"execution_time"`
	Weight            int32     `json:"weight" db:"weight"`
	Sets              string    `json:"sets" db:"sets"`
	Reps              string    `json:"reps" db:"reps"`
	RestTime          *int32    `json:"restTime,omitempty" db:"rest_time"`
	Thumbnail         *string   `json:"thumbnail,omitempty" db:"thumbnail"`
	TimeTotalWorkout  int32     `json:"timeTotalWorkout" db:"time_total_workout"`
	ExerciseTitle     string    `json:"exerciseTitle" db:"exercise_title"`
	ExerciseID        uuid.UUID `json:"exerciseId" db:"exercise_id"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
}

type Comment struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Content    string    `json:"content" db:"content"`
	StudentID  uuid.UUID `json:"studentId" db:"student_id"`
	PersonalID uuid.UUID `json:"personalId" db:"personal_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

type WorkoutsRating struct {
	ID         uuid.UUID `json:"id" db:"id"`
	StudentID  uuid.UUID `json:"studentId" db:"student_id"`
	WorkoutID  uuid.UUID `json:"workoutId" db:"workout_id"`
	PersonalID uuid.UUID `json:"personalId" db:"personal_id"`
	Rating     int32     `json:"rating" db:"rating"`
	Comment    *string   `json:"comment,omitempty" db:"comment"`
	RatingDate time.Time `json:"ratingDate" db:"rating_date"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

type PasswordResetToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"userId" db:"user_id"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expiresAt" db:"expires_at"`
	UsedAt    *time.Time `json:"usedAt,omitempty" db:"used_at"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
}

type RefreshToken struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"userId" db:"user_id"`
	Token      string     `json:"token" db:"token"`
	ExpiresAt  time.Time  `json:"expiresAt" db:"expires_at"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt" db:"updated_at"`
	RevokedAt  *time.Time `json:"revokedAt,omitempty" db:"revoked_at"`
	DeviceInfo *string    `json:"deviceInfo,omitempty" db:"device_info"`
	IPAddress  *string    `json:"ipAddress,omitempty" db:"ip_address"`
}

// Additional types for messages and other features
type CreateMessageRequest struct {
	PersonalID uuid.UUID `json:"personalId" validate:"required"`
	StudentID  uuid.UUID `json:"studentId" validate:"required"`
	Title      string    `json:"title" validate:"required,min=1,max=200"`
	Content    string    `json:"content" validate:"required,min=1,max=2000"`
}

type MessageResponse struct {
	ID         uuid.UUID `json:"id"`
	PersonalID uuid.UUID `json:"personalId"`
	StudentID  uuid.UUID `json:"studentId"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sentAt"`
}
// Request/Response types
type AuthenticateWithPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateStudentWithUserRequest struct {
	Name                  string    `json:"name" validate:"required,min=2,max=100"`
	Email                 string    `json:"email" validate:"required,email"`
	Phone                 string    `json:"phone" validate:"required"`
	Password              string    `json:"password" validate:"required,min=6"`
	BornDate              time.Time `json:"bornDate" validate:"required"`
	Age                   int32     `json:"age" validate:"required,min=13,max=120"`
	Weight                float64   `json:"weight" validate:"required,min=30,max=300"`
	Objective             string    `json:"objective" validate:"required"`
	TrainingFrequency     string    `json:"trainingFrequency" validate:"required"`
	DidBodybuilding       bool      `json:"didBodybuilding"`
	MedicalCondition      *string   `json:"medicalCondition,omitempty"`
	PhysicalActivityLevel *string   `json:"physicalActivityLevel,omitempty"`
	Observations          *string   `json:"observations,omitempty"`
}

type CreatePersonalWithUserRequest struct {
	Name           string  `json:"name" validate:"required,min=2,max=100"`
	Email          string  `json:"email" validate:"required,email"`
	Phone          string  `json:"phone" validate:"required"`
	Password       string  `json:"password" validate:"required,min=6"`
	Description    *string `json:"description,omitempty"`
	VideoURL       *string `json:"videoUrl,omitempty"`
	Experience     *string `json:"experience,omitempty"`
	Specialization *string `json:"specialization,omitempty"`
	Qualifications *string `json:"qualifications,omitempty"`
}

type UpdateProfileRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone *string `json:"phone,omitempty"`
}

type CreateWorkoutRequest struct {
	Name        string   `json:"name" validate:"required,min=2,max=100"`
	Description string   `json:"description" validate:"required"`
	Thumbnail   string   `json:"thumbnail"`
	WeekDays    *[]Day   `json:"weekDays,omitempty"`
	Exclusive   *bool    `json:"exclusive,omitempty"`
}

type UpdateWorkoutRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description,omitempty"`
	Thumbnail   *string `json:"thumbnail,omitempty"`
	WeekDays    *[]Day  `json:"weekDays,omitempty"`
	Exclusive   *bool   `json:"exclusive,omitempty"`
}

type CreateExerciseRequest struct {
	Name         string `json:"name" validate:"required,min=2,max=100"`
	Description  string `json:"description" validate:"required"`
	VideoURL     string `json:"videoUrl" validate:"required,url"`
	Instructions string `json:"instructions" validate:"required"`
	Category     string `json:"category" validate:"required"`
}

type UpdateExerciseRequest struct {
	Name         *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description  *string `json:"description,omitempty"`
	VideoURL     *string `json:"videoUrl,omitempty" validate:"omitempty,url"`
	Instructions *string `json:"instructions,omitempty"`
	Category     *string `json:"category,omitempty"`
}

type CreateSchedulingRequest struct {
	Date      time.Time      `json:"date" validate:"required"`
	StartTime time.Time      `json:"startTime" validate:"required"`
	EndTime   time.Time      `json:"endTime" validate:"required"`
	Type      SchedulingType `json:"type" validate:"required"`
	Notes     *string        `json:"notes,omitempty"`
	PersonalID uuid.UUID     `json:"personalId" validate:"required"`
}

type UpdateSchedulingRequest struct {
	Date      *time.Time      `json:"date,omitempty"`
	StartTime *time.Time      `json:"startTime,omitempty"`
	EndTime   *time.Time      `json:"endTime,omitempty"`
	Status    *SchedulingStatus `json:"status,omitempty"`
	Type      *SchedulingType `json:"type,omitempty"`
	Notes     *string         `json:"notes,omitempty"`
}
type WorkoutFrequencyResponse struct {
	TotalWorkouts   int32   `json:"total_workouts"`
	AveragePerWeek  float64 `json:"average_per_week"`
	MostActiveDay   string  `json:"most_active_day"`
	LongestStreak   int32   `json:"longest_streak"`
	CurrentStreak   int32   `json:"current_streak"`
	TotalDuration   int32   `json:"total_duration"`
	AverageDuration float64 `json:"average_duration"`
}

type WorkoutHistoryExerciseResponse struct {
	WorkoutHistoryID uuid.UUID `json:"workout_history_id"`
	WorkoutName      string    `json:"workout_name"`
	ExerciseName     string    `json:"exercise_name"`
	Sets             int32     `json:"sets"`
	Reps             int32     `json:"reps"`
	Weight           *float64  `json:"weight,omitempty"`
	Duration         int32     `json:"duration"`
	RestTime         *int32    `json:"rest_time,omitempty"`
	CompletedAt      time.Time `json:"completed_at"`
}

type ExerciseProgressPoint struct {
	Date   time.Time `json:"date"`
	Weight *float64  `json:"weight,omitempty"`
	Reps   int32     `json:"reps"`
	Volume float64   `json:"volume"`
}

type ExercisePerformanceResponse struct {
	ExerciseName      string                  `json:"exercise_name"`
	TotalSessions     int32                   `json:"total_sessions"`
	MaxWeight         *float64                `json:"max_weight,omitempty"`
	MaxReps           int32                   `json:"max_reps"`
	MaxVolume         float64                 `json:"max_volume"`
	AverageWeight     *float64                `json:"average_weight,omitempty"`
	AverageReps       float64                 `json:"average_reps"`
	AverageVolume     float64                 `json:"average_volume"`
	FirstRecorded     time.Time               `json:"first_recorded"`
	LastRecorded      time.Time               `json:"last_recorded"`
	ProgressData      []ExerciseProgressPoint `json:"progress_data"`
}

type UserStatisticsResponse struct {
	TotalWorkouts        int32      `json:"total_workouts"`
	TotalDuration        int32      `json:"total_duration"`
	AverageDuration      float64    `json:"average_duration"`
	FavoriteExercise     *string    `json:"favorite_exercise,omitempty"`
	TotalSchedulings     int32      `json:"total_schedulings"`
	CompletedSchedulings int32      `json:"completed_schedulings"`
	CancelledSchedulings int32      `json:"cancelled_schedulings"`
	JoinDate             time.Time  `json:"join_date"`
	LastActivity         *time.Time `json:"last_activity,omitempty"`
	CurrentStreak        int32      `json:"current_streak"`
	LongestStreak        int32      `json:"longest_streak"`
}

type PlatformStatisticsResponse struct {
	TotalUsers           int32   `json:"total_users"`
	TotalStudents        int32   `json:"total_students"`
	TotalTrainers        int32   `json:"total_trainers"`
	ActiveUsers          int32   `json:"active_users"`
	NewUsersThisMonth    int32   `json:"new_users_this_month"`
	TotalWorkouts        int32   `json:"total_workouts"`
	TotalExercises       int32   `json:"total_exercises"`
	WorkoutsThisMonth    int32   `json:"workouts_this_month"`
	TotalSchedulings     int32   `json:"total_schedulings"`
	CompletedSchedulings int32   `json:"completed_schedulings"`
	SchedulingsThisMonth int32   `json:"schedulings_this_month"`
	Revenue              float64 `json:"revenue"`
	RevenueThisMonth     float64 `json:"revenue_this_month"`
}

type DailySignupData struct {
	Date     time.Time `json:"date"`
	Students int32     `json:"students"`
	Trainers int32     `json:"trainers"`
	Total    int32     `json:"total"`
}

type UserReportResponse struct {
	TotalUsers      int32             `json:"total_users"`
	NewUsers        int32             `json:"new_users"`
	ActiveUsers     int32             `json:"active_users"`
	ChurnRate       float64           `json:"churn_rate"`
	RetentionRate   float64           `json:"retention_rate"`
	DailySignupData []DailySignupData `json:"daily_signup_data"`
}

type PopularExerciseData struct {
	ExerciseName string `json:"exercise_name"`
	UsageCount   int32  `json:"usage_count"`
	Category     string `json:"category"`
}

type WorkoutReportResponse struct {
	TotalWorkouts     int32                 `json:"total_workouts"`
	CompletedWorkouts int32                 `json:"completed_workouts"`
	AverageDuration   float64               `json:"average_duration"`
	TotalDuration     int32                 `json:"total_duration"`
	PopularExercises  []PopularExerciseData `json:"popular_exercises"`
}

type SchedulingReportResponse struct {
	TotalSchedulings     int32   `json:"total_schedulings"`
	CompletedSchedulings int32   `json:"completed_schedulings"`
	CancelledSchedulings int32   `json:"cancelled_schedulings"`
	PendingSchedulings   int32   `json:"pending_schedulings"`
	CompletionRate       float64 `json:"completion_rate"`
	CancellationRate     float64 `json:"cancellation_rate"`
	AverageLeadTime      float64 `json:"average_lead_time"`
}

type DailyRevenueData struct {
	Date    time.Time `json:"date"`
	Revenue float64   `json:"revenue"`
}

type RevenueReportResponse struct {
	TotalRevenue     float64            `json:"total_revenue"`
	AverageRevenue   float64            `json:"average_revenue"`
	RevenueGrowth    float64            `json:"revenue_growth"`
	DailyRevenueData []DailyRevenueData `json:"daily_revenue_data"`
}

type TopTrainerData struct {
	TrainerName       string  `json:"trainer_name"`
	TotalSchedulings  int32   `json:"total_schedulings"`
	CompletedSessions int32   `json:"completed_sessions"`
	Rating            float64 `json:"rating"`
	Revenue           float64 `json:"revenue"`
}

type TrainerReportResponse struct {
	TotalTrainers  int32            `json:"total_trainers"`
	ActiveTrainers int32            `json:"active_trainers"`
	AverageRating  float64          `json:"average_rating"`
	TotalSessions  int32            `json:"total_sessions"`
	TopTrainers    []TopTrainerData `json:"top_trainers"`
}

type WorkoutTrendPoint struct {
	Date            time.Time `json:"date"`
	WorkoutCount    int32     `json:"workout_count"`
	TotalDuration   int32     `json:"total_duration"`
	AverageDuration float64   `json:"average_duration"`
}

type WorkoutTrendsResponse struct {
	Period    string              `json:"period"`
	StartDate time.Time           `json:"start_date"`
	EndDate   time.Time           `json:"end_date"`
	TrendData []WorkoutTrendPoint `json:"trend_data"`
}

type MuscleGroupData struct {
	MuscleGroup   string     `json:"muscle_group"`
	WorkoutCount  int32      `json:"workout_count"`
	ExerciseCount int32      `json:"exercise_count"`
	TotalVolume   float64    `json:"total_volume"`
	LastWorked    *time.Time `json:"last_worked,omitempty"`
}

type MuscleGroupAnalysisResponse struct {
	MuscleGroups []MuscleGroupData `json:"muscle_groups"`
}
type PersonalTrainerResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	AvatarURL      *string   `json:"avatar_url,omitempty"`
	Description    *string   `json:"description,omitempty"`
	VideoURL       *string   `json:"video_url,omitempty"`
	Experience     *string   `json:"experience,omitempty"`
	Specialization *string   `json:"specialization,omitempty"`
	Qualifications *string   `json:"qualifications,omitempty"`
	Rating         *float64  `json:"rating,omitempty"`
	ReviewCount    int32     `json:"review_count"`
}

type TrainerCommentResponse struct {
	ID          uuid.UUID `json:"id"`
	StudentName string    `json:"student_name"`
	Comment     string    `json:"comment"`
	Rating      int32     `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
}

type StudentResponse struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Email                 string    `json:"email"`
	Phone                 string    `json:"phone"`
	AvatarURL             *string   `json:"avatar_url,omitempty"`
	BornDate              time.Time `json:"born_date"`
	Age                   int32     `json:"age"`
	Weight                float64   `json:"weight"`
	Objective             string    `json:"objective"`
	TrainingFrequency     string    `json:"training_frequency"`
	DidBodybuilding       bool      `json:"did_bodybuilding"`
	MedicalCondition      *string   `json:"medical_condition,omitempty"`
	PhysicalActivityLevel *string   `json:"physical_activity_level,omitempty"`
	Observations          *string   `json:"observations,omitempty"`
}

type SubscriptionStatus string

const (
	SubscriptionStatusPending   SubscriptionStatus = "PENDING"
	SubscriptionStatusActive    SubscriptionStatus = "ACTIVE"
	SubscriptionStatusCancelled SubscriptionStatus = "CANCELLED"
	SubscriptionStatusExpired   SubscriptionStatus = "EXPIRED"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
	UserStatusBanned    UserStatus = "BANNED"
)

type WorkoutHistoryResponse struct {
	ID          uuid.UUID `json:"id"`
	WorkoutName string    `json:"workout_name"`
	Duration    int32     `json:"duration"`
	Notes       *string   `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type ExerciseTemplateResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	VideoURL     string    `json:"video_url"`
	Instructions string    `json:"instructions"`
	Category     string    `json:"category"`
	MuscleGroups []string  `json:"muscle_groups"`
	Equipment    []string  `json:"equipment"`
	Difficulty   string    `json:"difficulty"`
}

type WorkoutTemplateResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	Category    string    `json:"category"`
	Difficulty  string    `json:"difficulty"`
	Duration    int32     `json:"duration"`
	WeekDays    []string  `json:"week_days"`
	Tags        []string  `json:"tags"`
}

type TrainingProgramResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	Category    string    `json:"category"`
	Difficulty  string    `json:"difficulty"`
	Duration    int32     `json:"duration"`
	IsFree      bool      `json:"is_free"`
	Price       *float64  `json:"price,omitempty"`
}

type AvailabilitySlot struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Available bool      `json:"available"`
}

type SchedulingHistoryResponse struct {
	ID            uuid.UUID        `json:"id"`
	Date          time.Time        `json:"date"`
	Type          SchedulingType   `json:"type"`
	Status        SchedulingStatus `json:"status"`
	WorkoutName   *string          `json:"workout_name,omitempty"`
	PartnerName   string           `json:"partner_name"`
	CompletedAt   *time.Time       `json:"completed_at,omitempty"`
}

type PlanResponse struct {
	ID              uuid.UUID `json:"id"`
	TrainerID       uuid.UUID `json:"trainer_id"`
	TrainerName     string    `json:"trainer_name,omitempty"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	Duration        int32     `json:"duration"`
	Features        []string  `json:"features"`
	IsActive        bool      `json:"is_active"`
	SubscriberCount int32     `json:"subscriber_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type SubscriptionResponse struct {
	ID        uuid.UUID          `json:"id"`
	UserID    uuid.UUID          `json:"user_id"`
	PlanID    uuid.UUID          `json:"plan_id"`
	PlanName  string             `json:"plan_name"`
	StartDate time.Time          `json:"start_date"`
	EndDate   time.Time          `json:"end_date"`
	Status    SubscriptionStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type SubscriptionHistoryResponse struct {
	ID          uuid.UUID          `json:"id"`
	PlanName    string             `json:"plan_name"`
	TrainerName string             `json:"trainer_name"`
	StartDate   time.Time          `json:"start_date"`
	EndDate     time.Time          `json:"end_date"`
	Status      SubscriptionStatus `json:"status"`
	Amount      float64            `json:"amount"`
	CreatedAt   time.Time          `json:"created_at"`
}

type SubscriberResponse struct {
	UserID       uuid.UUID          `json:"user_id"`
	UserName     string             `json:"user_name"`
	UserEmail    string             `json:"user_email"`
	StartDate    time.Time          `json:"start_date"`
	EndDate      time.Time          `json:"end_date"`
	Status       SubscriptionStatus `json:"status"`
	SubscribedAt time.Time          `json:"subscribed_at"`
}

type PlanRevenueResponse struct {
	PlanID            uuid.UUID `json:"plan_id"`
	PlanName          string    `json:"plan_name"`
	TotalRevenue      float64   `json:"total_revenue"`
	MonthlyRevenue    float64   `json:"monthly_revenue"`
	ActiveSubscribers int32     `json:"active_subscribers"`
	TotalSubscribers  int32     `json:"total_subscribers"`
	AverageRevenue    float64   `json:"average_revenue"`
}
