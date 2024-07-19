import Fastify from "fastify";
import * as grpc from "@grpc/grpc-js";
import {
  ProjectServiceService,
  UpdateProjectStatusRequest,
  UpdateProjectStatusResponse,
  ProjectStatus,
} from "./genprotobuf/project";
import { DBProjectStatus, updateStatusForProject } from "./repository/project";

const fastify = Fastify({ logger: true });

const server = new grpc.Server();

server.addService(ProjectServiceService, {
  updateProjectStatus: async (
    call: grpc.ServerUnaryCall<
      UpdateProjectStatusRequest,
      UpdateProjectStatusResponse
    >,
    callback: grpc.sendUnaryData<UpdateProjectStatusResponse>
  ) => {
    const { projectId, status } = call.request;
    console.log(
      `Updating project status: ${projectId} - ${ProjectStatus[status]}`
    );

    try {
      let statusToSet: DBProjectStatus;
      if (typeof status === "number") {
        switch (status) {
          case 0:
            statusToSet = "NOT_LIVE";
            break;
          case 1:
            statusToSet = "LIVE";
            break;
          case 2:
            statusToSet = "DEPLOYING";
            break;
          default:
            throw new Error("Invalid status value");
        }
      } else {
        statusToSet = status;
      }

      await updateStatusForProject(projectId, statusToSet);

      callback(null, {
        success: true,
        message: `Project ${projectId} status updated to ${ProjectStatus[status]} successfully`,
      });
    } catch (error) {
      console.error(error);
      callback(null, {
        success: false,
        message: `Failed to update project ${projectId} - expected status ${ProjectStatus[status]}`,
      });
    }
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
