-- Messages table
CREATE TABLE message (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    personal_id UUID NOT NULL REFERENCES personal(id),
    student_id UUID NOT NULL REFERENCES student(id),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

---- create above / drop below ----

-- Drop table
DROP TABLE IF EXISTS message;
