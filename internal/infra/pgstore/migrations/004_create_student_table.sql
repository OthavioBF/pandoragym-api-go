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

---- create above / drop below ----

-- Drop table
DROP TABLE IF EXISTS student;
