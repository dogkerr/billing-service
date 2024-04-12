CREATE TABLE charges_table (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  container_id TEXT NOT NULL,
  user_id VARCHAR(255) NOT NULL,
  total_cpu_usage FLOAT NOT NULL,
  total_memory_usage FLOAT NOT NULL,
  total_net_ingress_usage FLOAT NOT NULL,
  total_net_egress_usage FLOAT NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  total_cost FLOAT NOT NULL
);