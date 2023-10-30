-- SET statement_timeout = '5s';

CREATE TABLE IF NOT EXISTS orders (
  id SERIAL PRIMARY KEY NOT NULL,
  user_id INTEGER REFERENCES users(id) NOT NULL,
  amount INTEGER NOT NULL CHECK (amount >= 1),
  order_items TEXT NOT NULL,
  status VARCHAR(15) NOT NULL DEFAULT 'pending' CHECK (status IN ('paid', 'failed','pending','cancelled')),
  created_at timestamp WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
--bun:split
CREATE INDEX IF NOT EXISTS idx_order_user_id ON orders (user_id);
