-- Comments table
CREATE TABLE comment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    student_id UUID NOT NULL REFERENCES student(id),
    personal_id UUID NOT NULL REFERENCES personal(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

---- create above / drop below ----

-- Drop table
DROP TABLE IF EXISTS comment;
