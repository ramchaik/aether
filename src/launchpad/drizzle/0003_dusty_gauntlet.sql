DO $$ BEGIN
 CREATE TYPE "public"."project_status" AS ENUM('NOT_LIVE', 'LIVE', 'DEPLOYING');
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
ALTER TABLE "projects" ADD COLUMN "status" "project_status" DEFAULT 'NOT_LIVE';