import "dotenv/config";

import { clerkPlugin } from "@clerk/fastify";
import cors from "@fastify/cors";
import Fastify from "fastify";
import { pool } from "./db";
import { authMiddleware } from "./middlewares";
import registerRoutes from "./routes";

const PORT = (process.env.PORT || 8000) as number;
declare module "fastify" {
  interface FastifyRequest {
    userId: string;
  }
}
const app = Fastify({
  logger: true,
});

app.register(clerkPlugin);

app.register(cors, {
  origin: true,
  allowedHeaders: "*",
});

// auth
app.addHook("preHandler", authMiddleware);

app.register(registerRoutes);

app.get("/health", async (_req, reply) => {
  reply.code(200).send({ message: "Server is healthy!" });
});

app.addHook("onClose", async () => {
  await pool.end();
});

async function start() {
  try {
    await app.listen({
      host: "127.0.0.1", //"0.0.0.0",
      port: PORT,
    });
  } catch (error) {
    app.log.error(error);
    process.exit(1);
  }
}

start();
