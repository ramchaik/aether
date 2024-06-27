import Fastify from "fastify";
import registerRoutes from "./routes";
import prisma from "./lib/prisma";

const PORT = (process.env.PORT || 8000) as number;

const app = Fastify({
  logger: true,
});

app.register(registerRoutes);

app.get("/health", async (req, reply) => {
  reply.code(200).send({ message: "Server is healthy!" });
});

app.addHook("onClose", async () => {
  await prisma.$disconnect();
});

async function start() {
  try {
    await app.listen({ host: "0.0.0.0", port: PORT });
  } catch (error) {
    app.log.error(error);
    process.exit(1);
  }
}

start();
