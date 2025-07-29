-- Add ADMIN role to the role enum
-- ---- create ----
ALTER TYPE role ADD VALUE 'ADMIN';

-- ---- drop ----
-- Note: PostgreSQL doesn't support removing enum values directly
-- If rollback is needed, the entire enum would need to be recreated
-- For now, we'll leave this empty as removing enum values is complex
