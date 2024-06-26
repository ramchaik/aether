import { FastifyInstance } from "fastify";
import projectRoutes from "./project";

export default async function registerRoutes(fastify: FastifyInstance) {
  fastify.register(projectRoutes, { prefix: "/api/v1" });
}
