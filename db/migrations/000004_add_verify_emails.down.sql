DROP TABLE "verify_emails" IF EXISTS CASCADE;

ALTER TABLE "users" DROP COLUMN "is_email_verified";