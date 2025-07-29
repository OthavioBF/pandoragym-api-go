package pgstore

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateUserParams struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required,min=2,max=100"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Phone     string    `json:"phone" db:"phone" validate:"required"`
	Password  string    `json:"password" db:"password" validate:"required,min=6"`
	Role      Role      `json:"role" db:"role" validate:"required,oneof=PERSONAL STUDENT ADMIN"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateStudentParams struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	BornDate              time.Time `json:"bornDate" db:"born_date" validate:"required"`
	Age                   int32     `json:"age" db:"age" validate:"required,min=13,max=120"`
	Weight                float64   `json:"weight" db:"weight" validate:"required,min=30,max=300"`
	Objective             string    `json:"objective" db:"objective" validate:"required"`
	TrainingFrequency     string    `json:"trainingFrequency" db:"training_frequency" validate:"required"`
	DidBodybuilding       bool      `json:"didBodybuilding" db:"did_bodybuilding"`
	MedicalCondition      *string   `json:"medicalCondition,omitempty" db:"medical_condition"`
	PhysicalActivityLevel *string   `json:"physicalActivityLevel,omitempty" db:"physical_activity_level"`
	Observations          *string   `json:"observations,omitempty" db:"observations"`
}

type CreatePersonalParams struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Description    *string   `json:"description,omitempty" db:"description"`
	VideoURL       *string   `json:"videoUrl,omitempty" db:"video_url"`
	Experience     *string   `json:"experience,omitempty" db:"experience"`
	Specialization *string   `json:"specialization,omitempty" db:"specialization"`
	Qualifications *string   `json:"qualifications,omitempty" db:"qualifications"`
}

type UpdateUserParams struct {
	ID        uuid.UUID  `json:"id" db:"id" validate:"required"`
	Name      *string    `json:"name,omitempty" db:"name" validate:"omitempty,min=2,max=100"`
	Phone     *string    `json:"phone,omitempty" db:"phone" validate:"omitempty,min=10"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
}

type UpdateUserAvatarParams struct {
	ID        uuid.UUID `json:"id" db:"id" validate:"required"`
	AvatarURL string    `json:"avatarUrl" db:"avatar_url" validate:"required,url"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type UpdateUserPasswordParams struct {
	ID       uuid.UUID `json:"id" db:"id" validate:"required"`
	Password string    `json:"password" db:"password" validate:"required,min=6"`
}

type GetUserByEmailRow struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	AvatarURL *string   `json:"avatarUrl,omitempty" db:"avatar_url"`
	Password  string    `json:"-" db:"password"`
	Role      Role      `json:"role" db:"role"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type GetUserForAuthRow struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"-" db:"password"`
	Role     Role      `json:"role" db:"role"`
}

type GetUserByIDRow struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	AvatarURL *string   `json:"avatarUrl,omitempty" db:"avatar_url"`
	Role      Role      `json:"role" db:"role"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type GetUserByIdParams struct {
	ID uuid.UUID `json:"id" db:"id" validate:"required"`
}

type GetUserByIdRow struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	Name                  string     `json:"name" db:"name"`
	Email                 string     `json:"email" db:"email"`
	Phone                 string     `json:"phone" db:"phone"`
	AvatarURL             *string    `json:"avatarUrl,omitempty" db:"avatar_url"`
	Role                  Role       `json:"role" db:"role"`
	CreatedAt             time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time  `json:"updatedAt" db:"updated_at"`
	Rating                *float64   `json:"rating,omitempty" db:"rating"`
	Description           *string    `json:"description,omitempty" db:"description"`
	VideoURL              *string    `json:"videoUrl,omitempty" db:"video_url"`
	Experience            *string    `json:"experience,omitempty" db:"experience"`
	Specialization        *string    `json:"specialization,omitempty" db:"specialization"`
	Qualifications        *string    `json:"qualifications,omitempty" db:"qualifications"`
	BornDate              *time.Time `json:"bornDate,omitempty" db:"born_date"`
	Age                   *int32     `json:"age,omitempty" db:"age"`
	Weight                *float64   `json:"weight,omitempty" db:"weight"`
	Objective             *string    `json:"objective,omitempty" db:"objective"`
	TrainingFrequency     *string    `json:"trainingFrequency,omitempty" db:"training_frequency"`
	DidBodybuilding       *bool      `json:"didBodybuilding,omitempty" db:"did_bodybuilding"`
	MedicalCondition      *string    `json:"medicalCondition,omitempty" db:"medical_condition"`
	PhysicalActivityLevel *string    `json:"physicalActivityLevel,omitempty" db:"physical_activity_level"`
	Observations          *string    `json:"observations,omitempty" db:"observations"`
}

type GetAllUsersRow struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	AvatarURL *string   `json:"avatarUrl,omitempty" db:"avatar_url"`
	Role      Role      `json:"role" db:"role"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type UpdateUserProfileParams struct {
	ID        uuid.UUID `json:"id" db:"id" validate:"required"`
	Name      *string   `json:"name,omitempty" db:"name" validate:"omitempty,min=2,max=100"`
	Phone     *string   `json:"phone,omitempty" db:"phone" validate:"omitempty,min=10"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	AvatarURL *string   `json:"avatarUrl,omitempty"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Personal  *Personal `json:"personal,omitempty"`
	Student   *Student  `json:"student,omitempty"`
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  id, name, email, phone, password, role, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, name, email, phone, avatar_url, role, created_at, updated_at`

const createStudent = `-- name: CreateStudent :one
INSERT INTO student (
  id, born_date, age, weight, objective, training_frequency, did_bodybuilding, medical_condition, physical_activity_level, observations
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id, born_date, age, weight, objective, training_frequency, did_bodybuilding, medical_condition, physical_activity_level, observations`

const createPersonal = `-- name: CreatePersonal :one
INSERT INTO personal (
  id, description, video_url, experience, specialization, qualifications
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id, description, video_url, experience, specialization, qualifications`

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, phone, avatar_url, password, role, created_at, updated_at FROM users
WHERE email = $1 LIMIT 1`

const getUserForAuth = `-- name: GetUserForAuth :one
SELECT id, email, password, role FROM users
WHERE email = $1 LIMIT 1`

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, email, phone, avatar_url, role, created_at, updated_at FROM users
WHERE id = $1 LIMIT 1`

const updateUser = `-- name: UpdateUser :one
UPDATE users 
SET 
  name = COALESCE($2, name),
  phone = COALESCE($3, phone),
  updated_at = $4
WHERE id = $1
RETURNING id, name, email, phone, avatar_url, role, created_at, updated_at`

const updateUserAvatar = `-- name: UpdateUserAvatar :one
UPDATE users 
SET avatar_url = $2, updated_at = $3
WHERE id = $1
RETURNING id, name, email, phone, avatar_url, role, created_at, updated_at`

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users 
SET password = $2
WHERE id = $1`

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Phone,
		arg.Password,
		arg.Role,
		arg.CreatedAt,
		arg.UpdatedAt,
	)

	var i GetUserByEmailRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Phone,
		&i.AvatarURL,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}
	return i.ID, nil
}

func (q *Queries) CreateStudent(ctx context.Context, arg CreateStudentParams) error {
	_, err := q.db.Exec(ctx, createStudent,
		arg.ID,
		arg.BornDate,
		arg.Age,
		arg.Weight,
		arg.Objective,
		arg.TrainingFrequency,
		arg.DidBodybuilding,
		arg.MedicalCondition,
		arg.PhysicalActivityLevel,
		arg.Observations,
	)
	return err
}

func (q *Queries) CreatePersonal(ctx context.Context, arg CreatePersonalParams) error {
	_, err := q.db.Exec(ctx, createPersonal,
		arg.ID,
		arg.Description,
		arg.VideoURL,
		arg.Experience,
		arg.Specialization,
		arg.Qualifications,
	)
	return err
}

func (q *Queries) GetUserById(ctx context.Context, arg GetUserByIdParams) (*GetUserByIdRow, error) {
	row := q.db.QueryRow(ctx, getUserByID, arg.ID)

	var i GetUserByIdRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Phone,
		&i.AvatarURL,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*GetUserByEmailRow, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)

	var i GetUserByEmailRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Phone,
		&i.AvatarURL,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (q *Queries) GetUserForAuth(ctx context.Context, email string) (*GetUserForAuthRow, error) {
	row := q.db.QueryRow(ctx, getUserForAuth, email)

	var i GetUserForAuthRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Role,
	)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (q *Queries) GetAllUsers(ctx context.Context) ([]GetAllUsersRow, error) {
	// This would need a SQL query - placeholder for now
	return nil, nil
}

func (q *Queries) UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) error {
	_, err := q.db.Exec(ctx, updateUser, arg.ID, arg.Name, arg.Phone, arg.UpdatedAt)
	return err
}

func (q *Queries) UpdateUserAvatar(ctx context.Context, arg UpdateUserAvatarParams) error {
	_, err := q.db.Exec(ctx, updateUserAvatar, arg.ID, arg.AvatarURL, arg.UpdatedAt)
	return err
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.Exec(ctx, updateUserPassword, arg.ID, arg.Password)
	return err
}

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// This would need a SQL query - placeholder for now
	return nil
}

func (q *Queries) GetUserRole(ctx context.Context, id uuid.UUID) (Role, error) {
	// This would need a SQL query - placeholder for now
	return RoleStudent, nil
}

// Count queries for analytics

const countUsers = `-- name: CountUsers :one
SELECT COUNT(*) FROM users
`

func (q *Queries) CountUsers(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countUsers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countUsersByRole = `-- name: CountUsersByRole :one
SELECT COUNT(*) FROM users WHERE role = $1
`

func (q *Queries) CountUsersByRole(ctx context.Context, role Role) (int64, error) {
	row := q.db.QueryRow(ctx, countUsersByRole, role)
	var count int64
	err := row.Scan(&count)
	return count, err
}
