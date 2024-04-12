CREATE TABLE deposits_table (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id VARCHAR(255) NOT NULL,
  value FLOAT NOT NULL,
  status VARCHAR(255) NOT NULL
);