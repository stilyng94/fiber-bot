-- SET statement_timeout = '5s';

DROP INDEX IF EXISTS orders.idx_order_user_id;
--bun:split
DROP TABLE IF EXISTS orders;
