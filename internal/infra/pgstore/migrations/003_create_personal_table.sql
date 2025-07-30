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

---- create above / drop below ----

-- Drop table
DROP TABLE IF EXISTS personal;
