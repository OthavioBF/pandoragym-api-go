-- Plans table
CREATE TABLE plan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT[] DEFAULT '{}',
    price DECIMAL(10,2) NOT NULL,
    personal_id UUID REFERENCES personal(id)
);

---- create above / drop below ----

-- Drop table
DROP TABLE IF EXISTS plan;
