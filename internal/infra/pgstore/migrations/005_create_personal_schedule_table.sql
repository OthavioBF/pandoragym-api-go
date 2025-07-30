-- Personal schedules table
CREATE TABLE personal_schedule (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    week_day INTEGER NOT NULL,
    time_start_in_minutes INTEGER NOT NULL,
    time_end_in_minutes INTEGER NOT NULL,
    personal_id UUID NOT NULL REFERENCES personal(id) ON DELETE CASCADE
);

---- create above / drop below ----

-- Drop table
DROP TABLE IF EXISTS personal_schedule;
