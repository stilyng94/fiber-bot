-- SET statement_timeout = '5s';
CREATE TABLE IF NOT EXISTS products (
  id SERIAL PRIMARY KEY NOT NULL,
  title VARCHAR(75) NOT NULL,
  image_url TEXT,
  description VARCHAR(125) NOT NULL,
  price INTEGER NOT NULL CHECK (price >= 1),
  status VARCHAR(15) NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'unavailable')),
  created_at TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
--bun:split
CREATE UNIQUE INDEX IF NOT EXISTS uq_idx_products_title ON products (title);
