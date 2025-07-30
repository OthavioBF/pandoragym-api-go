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

---- create above / drop below ----

-- Drop table
DROP TABLE IF EXISTS exercises_template;
