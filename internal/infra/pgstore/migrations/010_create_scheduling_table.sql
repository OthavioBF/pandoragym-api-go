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

-- Create indexes for better performance
CREATE INDEX idx_scheduling_personal_id ON scheduling(personal_id);
CREATE INDEX idx_scheduling_student_id ON scheduling(student_id);
CREATE INDEX idx_scheduling_date ON scheduling(date);

---- create above / drop below ----

-- Drop indexes
DROP INDEX IF EXISTS idx_scheduling_date;
DROP INDEX IF EXISTS idx_scheduling_student_id;
DROP INDEX IF EXISTS idx_scheduling_personal_id;

-- Drop table
DROP TABLE IF EXISTS scheduling;
