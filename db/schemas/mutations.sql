CREATE TABLE mutations_table (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id VARCHAR(255) NOT NULL,
  mutation FLOAT NOT NULL,
  balance FLOAT NOT NULL,
  type ENUM('DEPOSIT', 'CHARGE') NOT NULL,
  deposit_id UUID REFERENCES deposits_table(id),
  charge_id UUID REFERENCES charges_table(id)
);