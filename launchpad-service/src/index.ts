import Fastify from "fastify";
import registerRoutes from "./routes";
import prisma from "./lib/prisma";

const PORT = (process.env.PORT || 3000) as number;

const fastify = Fastify({
  logger: true,
});

fastify.register(registerRoutes);

fastify.addHook("onClose", async () => {
  await prisma.$disconnect();
});

async function start() {
  try {
    await fastify.listen({ port: PORT });
  } catch (error) {
    fastify.log.error(error);
    process.exit(1);
  }
}

start();
