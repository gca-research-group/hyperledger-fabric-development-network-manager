CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL unique,
  password VARCHAR(255) NOT NULL,
  is_super BOOLEAN NOT NULL DEFAULT false,
  created_by INTEGER,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_by INTEGER,
  updated_at TIMESTAMP
);