-- Best practices:
-- 1. Do not change old migration files.
-- 2. Making your migrations idempotent.
-- 3. Use a clear and readable name.
-- Further notes: https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md
CREATE TABLE IF NOT EXISTS chat_history (
    thread_id VARCHAR(255) PRIMARY KEY,
    creator_id VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

