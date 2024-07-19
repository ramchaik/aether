import { FastifyInstance, FastifyReply, FastifyRequest } from "fastify";
// @ts-ignore
import { z } from "zod";
import { projectStatusEnum } from "../db/schema";
import * as repository from "../repository/project";
import { pushMessageToDeployQueue } from "../utils/awsSqs";
import { ERROR_MESSAGES, HTTP_CODES } from "../utils/httpCodes";

const createProjectSchema = z.object({
  name: z.string().min(1).max(100),
  customDomain: z.string().optional(),
  repositoryUrl: z.string(),
  buildCommand: z.string().optional(),
});

const PROXY_SVC = process.env.PROXY_SVC;

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
      // @ts-ignore
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

    // TODO: fix this when saving project
    // Update domain
    project.domain = `http://${project.id}.${PROXY_SVC}`;

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
    reply.log.error(error);
    reply.code(HTTP_CODES.INTERNAL_SERVER_ERROR).send({
      error: ERROR_MESSAGES.INTERNAL_SERVER_ERROR,
    });
  }
}

async function deployProjectHandler(
  request: FastifyRequest,
  reply: FastifyReply
) {
  try {
    const { id } = request.params as { id: string };
    const userId = request.userId;

    const project = await repository.readProject(userId, id);

    if (project.status === projectStatusEnum.enumValues[2]) {
      reply.code(HTTP_CODES.BAD_REQUEST).send({
        error: ERROR_MESSAGES.DEPLOYMENT_INPROGRESS,
      });
      return;
    }

    const message = {
      projectId: project.id,
      repoURL: project.repositoryUrl,
      buildCommand: project.buildCommand,
    };

    const messageId = await pushMessageToDeployQueue(message);

    // Update status a deploying
    await repository.updateStatusForProject(
      project.id,
      projectStatusEnum.enumValues[2]
    );

    reply.code(HTTP_CODES.OK).send({
      success: true,
      messageId: messageId,
    });
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

export default async function registerRoutes(fastify: FastifyInstance) {
  fastify.post<{ Body: CreateProjectBody }>("/project", createProjectHandler);
  fastify.post("/project/:id/deploy", deployProjectHandler);
  fastify.get("/project/:id", readProjectHandler);
  fastify.get("/project", readAllProjectHandler);
}
