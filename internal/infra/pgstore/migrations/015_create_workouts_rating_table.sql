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

---- create above / drop below ----

-- Drop table
DROP TABLE IF EXISTS workouts_rating;
