import { and, eq, sql } from "drizzle-orm";
import slugify from "slugify";
import { db } from "../db";
import { Project } from "../db/schema";
const { nanoid } = require("nanoid");

async function generateUniqueSlug(baseName: string): Promise<string> {
  const maxAttempts = 10;
  let attempts = 0;

  async function generateSlug() {
    const baseSlug = slugify(baseName, { lower: true, strict: true });
    const random = nanoid(6);
    const fullSlug = `${baseSlug}-${random}`;

    const existingSlug = await db
      .select()
      .from(Project)
      .where(eq(Project.slug, fullSlug));

    if (existingSlug && attempts < maxAttempts) {
      ++attempts;
      // Retry with a different random
      return generateSlug();
    }

    return fullSlug;
  }

  return generateSlug();
}

export async function checkDuplicateCustomDomain(
  domain: string
): Promise<boolean> {
  const existingCustomDomain = await db
    .select()
    .from(Project)
    .where(eq(Project.customDomain, domain));

  return existingCustomDomain.length !== 0;
}

interface ICreateProject {
  name: string;
  userId: string;
  repositoryUrl: string;
  customDomain?: string;
  buildCommand?: string;
}

export async function createProject({
  name,
  userId,
  repositoryUrl,
  customDomain,
  buildCommand,
}: ICreateProject) {
  const slug = await generateUniqueSlug(name);

  const domain = `${slug}.aether.app`;

  const projectData = {
    name: name,
    slug: slug,
    repositoryUrl: repositoryUrl,
    domain: domain,
    customDomain: customDomain,
    buildCommand: buildCommand,
    userId: userId,
  };

  const queryRes = await db
    .insert(Project)
    .values(projectData)
    .returning({ projectId: Project.id });
  return queryRes[0];
}

export async function readProject(userId: string, projectId: string) {
  const project = await db
    .select()
    .from(Project)
    .where(and(eq(Project.userId, userId), eq(Project.id, projectId)));

  if (project.length === 0) {
    throw new Error("Project: 404");
  }

  return project[0];
}

export async function readAllProjects(userId: string) {
  const projects = await db
    .select()
    .from(Project)
    .where(eq(Project.userId, userId));
  return projects;
}

export type DBProjectStatus = "NOT_LIVE" | "LIVE" | "DEPLOYING";
export async function updateStatusForProject(
  projectId: string,
  newStatus: DBProjectStatus
) {
  // First, check if the project exists and belongs to the user
  const existingProject = await db
    .select()
    .from(Project)
    .where(eq(Project.id, projectId));

  if (existingProject.length === 0) {
    throw new Error("Project not found");
  }

  // Update the project status
  const updatedProject = await db
    .update(Project)
    .set({ status: newStatus, updatedAt: sql`now()` })
    .where(eq(Project.id, projectId))
    .returning();

  if (updatedProject.length === 0) {
    throw new Error("Failed to update project status");
  }

  return updatedProject[0];
}
