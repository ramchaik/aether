import fastify, {
  FastifyInstance,
  FastifyReply,
  FastifyRequest,
} from "fastify";
// @ts-ignore
import { z } from "zod";
import { ERROR_MESSAGES, HTTP_CODES } from "../utils/httpCodes";
import * as repository from "../repository/project";

const createProjectSchema = z.object({
  name: z.string().min(1).max(100),
  customDomain: z.string().optional(),
  repositoryUrl: z.string(),
  buildCommand: z.string().optional(),
});

type CreateProjectBody = z.infer<typeof createProjectSchema>;

async function createProjectHandler(
  request: FastifyRequest,
  reply: FastifyReply
) {
  try {
    const {
      name,
      repositoryUrl,
      customDomain: _customDomain,
      buildCommand,
    } = createProjectSchema.parse(request.body);

    let customDomain = _customDomain;
    if (customDomain) {
      const isDuplicate = await repository.checkDuplicateCustomDomain(
        customDomain
      );
      if (isDuplicate) {
        return reply.code(HTTP_CODES.BAD_REQUEST).send({
          error: ERROR_MESSAGES.DUPLICATE_CUSTOM_DOMAIN,
        });
      }
    } else {
      customDomain = null;
    }

    const project = await repository.createProject({
      name,
      repositoryUrl,
      customDomain,
      buildCommand,
      userId: request.userId,
    });

    await reply.code(HTTP_CODES.CREATED).send(project);
  } catch (error) {
    if (error instanceof z.ZodError) {
      reply.code(HTTP_CODES.BAD_REQUEST).send({
        error: ERROR_MESSAGES.INVALID_INPUT,
        //@ts-ignore
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

async function readProjectHandler(
  request: FastifyRequest,
  reply: FastifyReply
) {
  try {
    const { id } = request.params as { id: string };
    const userId = request.userId;
    const project = await repository.readProject(userId, id);
    reply.code(HTTP_CODES.OK).send(project);
  } catch (error) {
    if ((error as Error).message === "Project: 404") {
      reply.code(HTTP_CODES.NOT_FOUND).send({
        error: ERROR_MESSAGES.PROJECT_NOT_FOUND,
      });
    } else {
      reply.log.error(error);
      reply.code(HTTP_CODES.INTERNAL_SERVER_ERROR).send({
        error: ERROR_MESSAGES.INTERNAL_SERVER_ERROR,
      });
    }
  }
}

async function readAllProjectHandler(
  request: FastifyRequest,
  reply: FastifyReply
) {
  try {
    const userId = request.userId;
    const projects = await repository.readAllProjects(userId);
    reply.code(HTTP_CODES.OK).send(projects);
  } catch (error) {
    reply.code(HTTP_CODES.INTERNAL_SERVER_ERROR).send({
      error: ERROR_MESSAGES.INTERNAL_SERVER_ERROR,
    });
  }
}

export default async function registerRoutes(fastify: FastifyInstance) {
  fastify.post<{ Body: CreateProjectBody }>("/project", createProjectHandler);
  fastify.get("/project/:id", readProjectHandler);
  fastify.get("/project", readAllProjectHandler);
}
