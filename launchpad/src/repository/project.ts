import slugify from "slugify";
import prisma from "../lib/prisma";

async function generateUniqueSlug(baseName: string): Promise<string> {
  const maxAttempts = 10;
  let attempts = 0;

  async function generateSlug() {
    const baseSlug = slugify(baseName, { lower: true, strict: true });
    const { nanoid } = await import("nanoid");
    const random = nanoid(6);
    const fullSlug = `${baseSlug}-${random}`;

    const existingSlug = await prisma.project.findUnique({
      where: { slug: fullSlug },
    });

    if (existingSlug && attempts < maxAttempts) {
      ++attempts;
      // Retry with a different random
      return generateSlug();
    }

    return fullSlug;
  }

  return generateSlug();
}

/**
 * @param domain custom domain
 * @returns true if duplicate otherwise false
 */
export async function checkDuplicateCustomDomain(
  domain: string
): Promise<boolean> {
  const exitingCustomDomain = await prisma.project.findUnique({
    where: { customDomain: domain },
  });

  return !!exitingCustomDomain;
}

interface ICreateProject {
  name: string;
  customDomain?: string;
  url: string;
}

export async function create({ name, customDomain, url }: ICreateProject) {
  const slug = await generateUniqueSlug(name);

  const domain = `${slug}.aether.app`;

  return prisma.project.create({
    data: {
      name,
      slug,
      gitURL: url,
      domain,
      customDomain,
    },
  });
}
