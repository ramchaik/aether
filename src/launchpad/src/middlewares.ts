import { getAuth } from "@clerk/fastify";
import { FastifyReply, FastifyRequest } from "fastify";

export const authMiddleware = async (
  request: FastifyRequest,
  reply: FastifyReply
) => {
  try {
    const auth = getAuth(request);
    if (!auth.userId) {
      return reply.code(401).send({ error: "Unauthorized" });
    }
    request.userId = auth.userId;
  } catch (error) {
    return reply.code(401).send({ error: "Invalid token" });
  }
};
