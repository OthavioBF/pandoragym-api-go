-- Scheduling history table
CREATE TABLE schedulings_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id UUID NOT NULL REFERENCES scheduling(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    status scheduling_status NOT NULL,
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    changed_by VARCHAR(255) NOT NULL,
    reason TEXT,
    notes TEXT
);

---- create above / drop below ----

-- Drop table
DROP TABLE IF EXISTS schedulings_history;
