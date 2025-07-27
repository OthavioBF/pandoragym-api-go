package pgstore

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// User-related request/response types
type CreateUserParams struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required,min=2,max=100"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Phone     string    `json:"phone" db:"phone" validate:"required"`
	Password  string    `json:"password" db:"password" validate:"required,min=6"`
	Role      Role      `json:"role" db:"role" validate:"required,oneof=PERSONAL STUDENT"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateStudentParams struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	BornDate              time.Time `json:"born_date" db:"born_date" validate:"required"`
	Age                   int32     `json:"age" db:"age" validate:"required,min=13,max=120"`
	Weight                float64   `json:"weight" db:"weight" validate:"required,min=30,max=300"`
	Objective             string    `json:"objective" db:"objective" validate:"required"`
	TrainingFrequency     string    `json:"training_frequency" db:"training_frequency" validate:"required"`
	DidBodybuilding       bool      `json:"did_bodybuilding" db:"did_bodybuilding"`
	MedicalCondition      *string   `json:"medical_condition,omitempty" db:"medical_condition"`
	PhysicalActivityLevel *string   `json:"physical_activity_level,omitempty" db:"physical_activity_level"`
	Observations          *string   `json:"observations,omitempty" db:"observations"`
}

type CreatePersonalParams struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Description    *string   `json:"description,omitempty" db:"description"`
	VideoURL       *string   `json:"video_url,omitempty" db:"video_url"`
	Experience     *string   `json:"experience,omitempty" db:"experience"`
	Specialization *string   `json:"specialization,omitempty" db:"specialization"`
	Qualifications *string   `json:"qualifications,omitempty" db:"qualifications"`
}

type CreateStudentWithUserRequest struct {
	Name                  string    `json:"name" validate:"required,min=2,max=100"`
	Email                 string    `json:"email" validate:"required,email"`
	Phone                 string    `json:"phone" validate:"required"`
	Password              string    `json:"password" validate:"required,min=6"`
	BornDate              time.Time `json:"born_date" validate:"required"`
	Age                   int32     `json:"age" validate:"required,min=13,max=120"`
	Weight                float64   `json:"weight" validate:"required,min=30,max=300"`
	Objective             string    `json:"objective" validate:"required"`
	TrainingFrequency     string    `json:"training_frequency" validate:"required"`
	DidBodybuilding       bool      `json:"did_bodybuilding"`
	MedicalCondition      *string   `json:"medical_condition,omitempty"`
	PhysicalActivityLevel *string   `json:"physical_activity_level,omitempty"`
	Observations          *string   `json:"observations,omitempty"`
}

type CreatePersonalWithUserRequest struct {
	Name           string  `json:"name" validate:"required,min=2,max=100"`
	Email          string  `json:"email" validate:"required,email"`
	Phone          string  `json:"phone" validate:"required"`
	Password       string  `json:"password" validate:"required,min=6"`
	Description    *string `json:"description,omitempty"`
	PresentationVideo *string `json:"presentation_video,omitempty"`
	Experience     *string `json:"experience,omitempty"`
	Specialization *string `json:"specialization,omitempty"`
	Qualifications *string `json:"qualifications,omitempty"`
}

type GetUserByIdParams struct {
	ID uuid.UUID `json:"id" db:"id" validate:"required"`
}

type UpdateUserProfileParams struct {
	ID        uuid.UUID  `json:"id" db:"id" validate:"required"`
	Name      *string    `json:"name,omitempty" db:"name" validate:"omitempty,min=2,max=100"`
	Phone     *string    `json:"phone,omitempty" db:"phone" validate:"omitempty,min=10"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type UpdateUserAvatarParams struct {
	ID        uuid.UUID `json:"id" db:"id" validate:"required"`
	AvatarURL string    `json:"avatar_url" db:"avatar_url" validate:"required,url"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UpdateUserPasswordParams struct {
	ID       uuid.UUID `json:"id" db:"id" validate:"required"`
	Password string    `json:"password" db:"password" validate:"required,min=6"`
}

type UpdateProfileRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone *string `json:"phone,omitempty" validate:"omitempty,min=10"`
}

type GetUserByIdRow struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	Name                  string     `json:"name" db:"name"`
	Email                 string     `json:"email" db:"email"`
	Phone                 string     `json:"phone" db:"phone"`
	AvatarURL             *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	Role                  Role       `json:"role" db:"role"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
	Rating                *float64   `json:"rating,omitempty" db:"rating"`
	Description           *string    `json:"description,omitempty" db:"description"`
	VideoURL              *string    `json:"video_url,omitempty" db:"video_url"`
	Experience            *string    `json:"experience,omitempty" db:"experience"`
	Specialization        *string    `json:"specialization,omitempty" db:"specialization"`
	Qualifications        *string    `json:"qualifications,omitempty" db:"qualifications"`
	BornDate              *time.Time `json:"born_date,omitempty" db:"born_date"`
	Age                   *int32     `json:"age,omitempty" db:"age"`
	Weight                *float64   `json:"weight,omitempty" db:"weight"`
	Objective             *string    `json:"objective,omitempty" db:"objective"`
	TrainingFrequency     *string    `json:"training_frequency,omitempty" db:"training_frequency"`
	DidBodybuilding       *bool      `json:"did_bodybuilding,omitempty" db:"did_bodybuilding"`
	MedicalCondition      *string    `json:"medical_condition,omitempty" db:"medical_condition"`
	PhysicalActivityLevel *string    `json:"physical_activity_level,omitempty" db:"physical_activity_level"`
	Observations          *string    `json:"observations,omitempty" db:"observations"`
}

type GetUserByEmailRow struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	AvatarURL *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	Password  string    `json:"-" db:"password"`
	Role      Role      `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type GetUserForAuthRow struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"-" db:"password"`
	Role     Role      `json:"role" db:"role"`
}

type GetAllUsersRow struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	AvatarURL *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	Role      Role      `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Personal  *Personal `json:"personal,omitempty"`
	Student   *Student  `json:"student,omitempty"`
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, name, email, phone, password, role, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id`

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
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createStudent = `-- name: CreateStudent :exec
INSERT INTO student (id, born_date, age, weight, objective, training_frequency, did_bodybuilding, medical_condition, physical_activity_level, observations)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

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

const createPersonal = `-- name: CreatePersonal :exec
INSERT INTO personal (id, description, video_url, experience, specialization, qualifications)
VALUES ($1, $2, $3, $4, $5, $6)`

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

const getUserById = `-- name: GetUserById :one
SELECT u.id, u.name, u.email, u.phone, u.avatar_url, u.role, u.created_at, u.updated_at,
       p.rating, p.description, p.video_url, p.experience, p.specialization, p.qualifications,
       st.born_date, st.age, st.weight, st.objective, st.training_frequency, st.did_bodybuilding,
       st.medical_condition, st.physical_activity_level, st.observations
FROM users u
LEFT JOIN personal p ON u.id = p.id
LEFT JOIN student st ON u.id = st.id
WHERE u.id = $1`

func (q *Queries) GetUserById(ctx context.Context, arg GetUserByIdParams) (*GetUserByIdRow, error) {
	row := q.db.QueryRow(ctx, getUserById, arg.ID)
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
		&i.Rating,
		&i.Description,
		&i.VideoURL,
		&i.Experience,
		&i.Specialization,
		&i.Qualifications,
		&i.BornDate,
		&i.Age,
		&i.Weight,
		&i.Objective,
		&i.TrainingFrequency,
		&i.DidBodybuilding,
		&i.MedicalCondition,
		&i.PhysicalActivityLevel,
		&i.Observations,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, phone, avatar_url, password, role, created_at, updated_at
FROM users
WHERE email = $1`

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
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}

const getUserForAuth = `-- name: GetUserForAuth :one
SELECT id, email, password, role FROM users WHERE email = $1`

func (q *Queries) GetUserForAuth(ctx context.Context, email string) (*GetUserForAuthRow, error) {
	row := q.db.QueryRow(ctx, getUserForAuth, email)
	var i GetUserForAuthRow
	err := row.Scan(&i.ID, &i.Email, &i.Password, &i.Role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT id, name, email, phone, avatar_url, role, created_at, updated_at
FROM users
ORDER BY created_at DESC`

func (q *Queries) GetAllUsers(ctx context.Context) ([]GetAllUsersRow, error) {
	rows, err := q.db.Query(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllUsersRow
	for rows.Next() {
		var i GetAllUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Phone,
			&i.AvatarURL,
			&i.Role,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserProfile = `-- name: UpdateUserProfile :exec
UPDATE users 
SET name = COALESCE($2, name), 
    phone = COALESCE($3, phone), 
    updated_at = $4
WHERE id = $1`

func (q *Queries) UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) error {
	_, err := q.db.Exec(ctx, updateUserProfile,
		arg.ID,
		arg.Name,
		arg.Phone,
		arg.UpdatedAt,
	)
	return err
}

const updateUserAvatar = `-- name: UpdateUserAvatar :exec
UPDATE users 
SET avatar_url = $2, updated_at = $3
WHERE id = $1`

func (q *Queries) UpdateUserAvatar(ctx context.Context, arg UpdateUserAvatarParams) error {
	_, err := q.db.Exec(ctx, updateUserAvatar, arg.ID, arg.AvatarURL, arg.UpdatedAt)
	return err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users SET password = $2 WHERE id = $1`

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.Exec(ctx, updateUserPassword, arg.ID, arg.Password)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUserRole = `-- name: GetUserRole :one
SELECT role FROM users WHERE id = $1`

func (q *Queries) GetUserRole(ctx context.Context, id uuid.UUID) (Role, error) {
	row := q.db.QueryRow(ctx, getUserRole, id)
	var role Role
	err := row.Scan(&role)
	return role, err
}
