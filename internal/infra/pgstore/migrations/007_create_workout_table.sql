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

-- Create indexes for better performance
CREATE INDEX idx_workout_personal_id ON workout(personal_id);
CREATE INDEX idx_workout_student_id ON workout(student_id);

---- create above / drop below ----

-- Drop indexes
DROP INDEX IF EXISTS idx_workout_student_id;
DROP INDEX IF EXISTS idx_workout_personal_id;

-- Drop table
DROP TABLE IF EXISTS workout;
