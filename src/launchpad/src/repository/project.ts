import { db } from "../db";
import slugify from "slugify";
import { Project } from "../db/schema";
import { eq } from "drizzle-orm";
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
  customDomain?: string;
  repositoryUrl: string;
}

export async function createProject({
  name,
  customDomain,
  repositoryUrl,
}: ICreateProject) {
  const slug = await generateUniqueSlug(name);

  const domain = `${slug}.aether.app`;

  const projectData = {
    name: name,
    slug: slug,
    repositoryUrl: repositoryUrl,
    domain: domain,
    customDomain: customDomain,
  };

  const queryRes = await db
    .insert(Project)
    .values(projectData)
    .returning({ projectId: Project.id });
  return queryRes[0];
}

export async function readProject(projectId: string) {
  const project = await db
    .select()
    .from(Project)
    .where(eq(Project.id, projectId));

  if (project.length === 0) {
    throw new Error("Project: 404");
  }

  return project[0];
}
