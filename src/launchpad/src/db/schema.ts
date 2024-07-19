import { pgTable, uuid, varchar, timestamp, pgEnum } from "drizzle-orm/pg-core";
import { sql } from "drizzle-orm";

export const projectStatusEnum = pgEnum("project_status", [
  "NOT_LIVE",
  "LIVE",
  "DEPLOYING",
]);

export const Project = pgTable("projects", {
  id: uuid("id")
    .default(sql`gen_random_uuid()`)
    .primaryKey(),
  name: varchar("name"),
  slug: varchar("slug").unique(),
  domain: varchar("domain")
    .unique()
    .default(sql`null`),
  repositoryUrl: varchar("repository_url"),
  customDomain: varchar("custom_domain")
    .unique()
    .default(sql`null`),
  buildCommand: varchar("build_command").default("npm run build"),
  createdAt: timestamp("created_at").default(sql`now()`),
  updatedAt: timestamp("updated_at")
    .default(sql`now()`)
    .$onUpdateFn(() => sql`now()`),
  userId: varchar("clerk_user_id"),
  status: projectStatusEnum("status").default("NOT_LIVE"),
});
