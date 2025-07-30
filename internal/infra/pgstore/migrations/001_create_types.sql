-- Create custom types
CREATE TYPE role AS ENUM ('PERSONAL', 'STUDENT', 'ADMIN');
CREATE TYPE scheduling_status AS ENUM ('PENDING_CONFIRMATION', 'SCHEDULED', 'IN_PROGRESS', 'RESCHEDULED', 'CANCELED', 'COMPLETED', 'MISSED');
CREATE TYPE scheduling_type AS ENUM ('ONLINE', 'IN_PERSON');
CREATE TYPE level AS ENUM ('BEGINNER', 'INTERMEDIARY', 'ADVANCED');
CREATE TYPE day AS ENUM ('Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sab');

---- create above / drop below ----

-- Drop custom types
DROP TYPE IF EXISTS day;
DROP TYPE IF EXISTS level;
DROP TYPE IF EXISTS scheduling_type;
DROP TYPE IF EXISTS scheduling_status;
DROP TYPE IF EXISTS role;
