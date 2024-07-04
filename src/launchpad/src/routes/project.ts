import fastify, {
  FastifyInstance,
  FastifyReply,
  FastifyRequest,
} from "fastify";
import { z } from "zod";
import { ERROR_MESSAGES, HTTP_CODES } from "../utils/httpCodes";
import * as repository from "../repository/project";

const createProjectSchema = z.object({
  name: z.string().min(1).max(100),
  customDomain: z.string().optional(),
  url: z.string(),
});

type CreateProjectBody = z.infer<typeof createProjectSchema>;

async function createProjectHandler(
  request: FastifyRequest,
  reply: FastifyReply
) {
  try {
    const { name, url, customDomain } = createProjectSchema.parse(request.body);

    if (customDomain) {
      const isDuplicate = await repository.checkDuplicateCustomDomain(
        customDomain
      );
      if (isDuplicate) {
        return reply.code(HTTP_CODES.BAD_REQUEST).send({
          error: ERROR_MESSAGES.DUPLICATE_CUSTOM_DOMAIN,
        });
      }
    }

    const project = await repository.create({
      name,
      url,
      customDomain,
    });

    await reply.code(HTTP_CODES.CREATED).send(project);
  } catch (error) {
    if (error instanceof z.ZodError) {
      reply.code(HTTP_CODES.BAD_REQUEST).send({
        error: ERROR_MESSAGES.INVALID_INPUT,
        details: error.errors,
      });
    } else {
      reply.log.error(error);
      reply.code(HTTP_CODES.INTERNAL_SERVER_ERROR).send({
        error: ERROR_MESSAGES.INTERNAL_SERVER_ERROR,
      });
    }
  }
}

export default async function registerRoutes(fastify: FastifyInstance) {
  fastify.post<{ Body: CreateProjectBody }>("/project", createProjectHandler);
}
