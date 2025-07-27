-- Create custom types
CREATE TYPE role AS ENUM ('PERSONAL', 'STUDENT');
CREATE TYPE scheduling_status AS ENUM ('PENDING_CONFIRMATION', 'SCHEDULED', 'IN_PROGRESS', 'RESCHEDULED', 'CANCELED', 'COMPLETED', 'MISSED');
CREATE TYPE scheduling_type AS ENUM ('ONLINE', 'IN_PERSON');
CREATE TYPE level AS ENUM ('BEGINNER', 'INTERMEDIARY', 'ADVANCED');
CREATE TYPE day AS ENUM ('Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sab');

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) NOT NULL,
    avatar_url TEXT,
    password VARCHAR(255) NOT NULL,
    role role NOT NULL,
    student_id UUID NOT NULL DEFAULT gen_random_uuid(),
    personal_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Personal trainers table
CREATE TABLE personal (
    id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    plan_id TEXT,
    rating DECIMAL(3,2),
    description TEXT,
    video_url TEXT,
    experience TEXT,
    specialization TEXT,
    qualifications TEXT
);

-- Students table
CREATE TABLE student (
    id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    born_date DATE NOT NULL,
    age INTEGER NOT NULL,
    weight DECIMAL(5,2) NOT NULL,
    objective TEXT NOT NULL,
    training_frequency TEXT NOT NULL,
    did_bodybuilding BOOLEAN DEFAULT FALSE,
    medical_condition TEXT,
    physical_activity_level TEXT,
    observations TEXT,
    personal_id UUID REFERENCES personal(id),
    plan_id TEXT
);

-- Personal schedules table
CREATE TABLE personal_schedule (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    week_day INTEGER NOT NULL,
    time_start_in_minutes INTEGER NOT NULL,
    time_end_in_minutes INTEGER NOT NULL,
    personal_id UUID NOT NULL REFERENCES personal(id) ON DELETE CASCADE
);

-- Plans table
CREATE TABLE plan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT[] DEFAULT '{}',
    price DECIMAL(10,2) NOT NULL,
    personal_id UUID REFERENCES personal(id)
);

-- Workouts table
CREATE TABLE workout (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    thumbnail TEXT NOT NULL,
    video_url TEXT,
    rest_time_between_exercises INTEGER,
    level level,
    week_days day[] DEFAULT '{}',
    exclusive BOOLEAN DEFAULT FALSE,
    is_template BOOLEAN DEFAULT FALSE,
    modality VARCHAR(255) NOT NULL,
    personal_id UUID REFERENCES personal(id),
    student_id UUID REFERENCES student(id),
    plan_id TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Exercises template table
CREATE TABLE exercises_template (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    thumbnail TEXT NOT NULL,
    video_url TEXT NOT NULL,
    load INTEGER,
    sets INTEGER NOT NULL,
    reps INTEGER NOT NULL,
    rest_time_between_sets INTEGER,
    personal_id UUID REFERENCES personal(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Exercises setup table (exercises within workouts)
CREATE TABLE exercises_setup (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    thumbnail TEXT NOT NULL,
    video_url TEXT NOT NULL,
    sets INTEGER NOT NULL,
    reps INTEGER NOT NULL,
    rest_time_between_sets INTEGER NOT NULL,
    load INTEGER NOT NULL,
    workout_id UUID NOT NULL REFERENCES workout(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Scheduling table
CREATE TABLE scheduling (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    personal_id UUID NOT NULL REFERENCES personal(id),
    student_id UUID NOT NULL REFERENCES student(id),
    workout_id UUID REFERENCES workout(id),
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    type scheduling_type NOT NULL,
    status scheduling_status NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    stard_at TIMESTAMP WITH TIME ZONE, -- Note: keeping original typo for compatibility
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    user_id UUID REFERENCES users(id)
);

-- Scheduling history table
CREATE TABLE schedulings_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id UUID NOT NULL REFERENCES scheduling(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    status scheduling_status NOT NULL,
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    changed_by VARCHAR(255) NOT NULL,
    reason TEXT,
    notes TEXT
);

-- Messages table
CREATE TABLE message (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    personal_id UUID NOT NULL REFERENCES personal(id),
    student_id UUID NOT NULL REFERENCES student(id),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Workouts history table
CREATE TABLE workouts_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES student(id),
    workout_id UUID NOT NULL REFERENCES workout(id),
    execution_time TEXT,
    weight INTEGER NOT NULL,
    sets TEXT NOT NULL,
    reps TEXT NOT NULL,
    rest_time INTEGER,
    thumbnail TEXT,
    time_total_workout INTEGER NOT NULL,
    exercise_title VARCHAR(255) NOT NULL,
    exercise_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Comments table
CREATE TABLE comment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    student_id UUID NOT NULL REFERENCES student(id),
    personal_id UUID NOT NULL REFERENCES personal(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Workouts rating table
CREATE TABLE workouts_rating (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES student(id),
    workout_id UUID NOT NULL REFERENCES workout(id),
    personal_id UUID NOT NULL REFERENCES personal(id),
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    rating_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Password reset tokens table
CREATE TABLE password_reset_token (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_workout_personal_id ON workout(personal_id);
CREATE INDEX idx_workout_student_id ON workout(student_id);
CREATE INDEX idx_scheduling_personal_id ON scheduling(personal_id);
CREATE INDEX idx_scheduling_student_id ON scheduling(student_id);
CREATE INDEX idx_scheduling_date ON scheduling(date);
CREATE INDEX idx_exercises_setup_workout_id ON exercises_setup(workout_id);
CREATE INDEX idx_password_reset_token_token ON password_reset_token(token);
CREATE INDEX idx_password_reset_token_user_id ON password_reset_token(user_id);

---- create above / drop below ----

-- Drop indexes
DROP INDEX IF EXISTS idx_password_reset_token_user_id;
DROP INDEX IF EXISTS idx_password_reset_token_token;
DROP INDEX IF EXISTS idx_exercises_setup_workout_id;
DROP INDEX IF EXISTS idx_scheduling_date;
DROP INDEX IF EXISTS idx_scheduling_student_id;
DROP INDEX IF EXISTS idx_scheduling_personal_id;
DROP INDEX IF EXISTS idx_workout_student_id;
DROP INDEX IF EXISTS idx_workout_personal_id;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS password_reset_token;
DROP TABLE IF EXISTS workouts_rating;
DROP TABLE IF EXISTS comment;
DROP TABLE IF EXISTS workouts_history;
DROP TABLE IF EXISTS message;
DROP TABLE IF EXISTS schedulings_history;
DROP TABLE IF EXISTS scheduling;
DROP TABLE IF EXISTS exercises_setup;
DROP TABLE IF EXISTS exercises_template;
DROP TABLE IF EXISTS workout;
DROP TABLE IF EXISTS plan;
DROP TABLE IF EXISTS personal_schedule;
DROP TABLE IF EXISTS student;
DROP TABLE IF EXISTS personal;
DROP TABLE IF EXISTS users;

-- Drop custom types
DROP TYPE IF EXISTS day;
DROP TYPE IF EXISTS level;
DROP TYPE IF EXISTS scheduling_type;
DROP TYPE IF EXISTS scheduling_status;
DROP TYPE IF EXISTS role;
