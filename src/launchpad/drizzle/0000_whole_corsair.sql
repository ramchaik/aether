CREATE TABLE IF NOT EXISTS "projects" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"name" varchar,
	"slug" varchar,
	"domain" varchar DEFAULT null,
	"repository_url" varchar,
	"custom_domain" varchar DEFAULT null,
	"created_at" timestamp DEFAULT now(),
	"updated_at" timestamp DEFAULT now(),
	CONSTRAINT "projects_slug_unique" UNIQUE("slug"),
	CONSTRAINT "projects_domain_unique" UNIQUE("domain"),
	CONSTRAINT "projects_custom_domain_unique" UNIQUE("custom_domain")
);
