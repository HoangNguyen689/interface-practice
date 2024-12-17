-- Best practices:
-- 1. Do not change old migration files.
-- 2. Making your migrations idempotent.
-- 3. Use a clear and readable name.
-- Further notes: https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md
DROP TABLE IF EXISTS chat_history_message;