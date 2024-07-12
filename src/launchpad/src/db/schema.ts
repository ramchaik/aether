import { pgTable, uuid, varchar, timestamp } from "drizzle-orm/pg-core";
import { sql } from "drizzle-orm";

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
  createdAt: timestamp("created_at").default(sql`now()`),
  updatedAt: timestamp("updated_at")
    .default(sql`now()`)
    .$onUpdateFn(() => sql`now()`),
});
