CREATE TABLE IF NOT EXISTS users  (
  id SERIAL PRIMARY KEY NOT NULL,
  telegram_id INTEGER NOT NULL,
  first_name VARCHAR(75),
  last_name VARCHAR(75),
  user_name VARCHAR(75),
  chat_id INTEGER,
  chat_type VARCHAR(75),
  role VARCHAR(15) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
  created_at TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS uq_idx_users_telegram_id ON users (telegram_id);
