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

---- create above / drop below ----

-- Drop table
DROP TABLE IF EXISTS workouts_history;
