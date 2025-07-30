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

-- Create index for better performance
CREATE INDEX idx_exercises_setup_workout_id ON exercises_setup(workout_id);

---- create above / drop below ----

-- Drop index
DROP INDEX IF EXISTS idx_exercises_setup_workout_id;

-- Drop table
DROP TABLE IF EXISTS exercises_setup;
