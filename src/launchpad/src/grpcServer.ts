import Fastify from "fastify";
import * as grpc from "@grpc/grpc-js";
import {
  ProjectServiceService,
  SaveProjectUrlRequest,
  SaveProjectUrlResponse,
} from "./proto/project";

const fastify = Fastify({ logger: true });

const server = new grpc.Server();

server.addService(ProjectServiceService, {
  saveProjectUrl: (
    call: grpc.ServerUnaryCall<SaveProjectUrlRequest, SaveProjectUrlResponse>,
    callback: grpc.sendUnaryData<SaveProjectUrlResponse>
  ) => {
    const { projectUrl, projectId } = call.request;
    console.log(`Saving project: ${projectId} - ${projectUrl}`);

    callback(null, { success: true, message: "Project saved successfully" });
  },
});

export const startGrpcServer = async () => {
  const GRPC_SERVER_ADDRESS = process.env.GRPC_PORT ?? "0.0.0.0:50051";
  await fastify.ready();
  server.bindAsync(
    GRPC_SERVER_ADDRESS,
    grpc.ServerCredentials.createInsecure(),
    (err, port) => {
      if (err) {
        console.error("Failed to bind server:", err);
        return;
      }
      fastify.log.info(`GRPC Server running at http://0.0.0.0:${port}`);
      server.start();
    }
  );
};
