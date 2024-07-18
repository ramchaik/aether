"use client";

import React from "react";
import {
  Card,
  CardBody,
  CardHeader,
  Divider,
  Chip,
  Button,
  Tooltip,
} from "@nextui-org/react";
import { motion } from "framer-motion";
import { useParams } from "next/navigation";
import { useFetchProject } from "@/hooks/useProjectApi";
import { Project } from "@/store/useProjectStore";

// Dummy build logs
const buildLogs = [
  { timestamp: "2024-07-11T10:05:00Z", message: "Build started" },
  { timestamp: "2024-07-11T10:05:05Z", message: "Installing dependencies..." },
  {
    timestamp: "2024-07-11T10:06:00Z",
    message: "Dependencies installed successfully",
  },
  {
    timestamp: "2024-07-11T10:06:05Z",
    message: "Running build command: npm run build",
  },
  {
    timestamp: "2024-07-11T10:07:00Z",
    message: "Build completed successfully",
  },
  { timestamp: "2024-07-11T10:07:05Z", message: "Deploying to Vercel..." },
  { timestamp: "2024-07-11T10:08:00Z", message: "Deployment successful" },
];

const ProjectDetailPage: React.FC = () => {
  const params = useParams();
  const projectId = params.id as string;
  const {
    data: project,
    isLoading,
    // error,
  } = useFetchProject<Project>(projectId);

  if (isLoading) return <div>Loading...</div>;
  if (!project) return <div>Project not found</div>;
  return (
    <div className="container mx-auto px-4 py-10">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <Card className="mb-6">
          <CardHeader className="flex justify-between items-center">
            <h1 className="text-2xl font-bold">{project.name}</h1>
            {/* <Chip color={project.status === "Live" ? "success" : "warning"}>
              {project.status}
            </Chip> */}
          </CardHeader>
          <Divider />
          <CardBody>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <p>
                  <span className="font-semibold">Project ID:</span>{" "}
                  {project.id}
                </p>
                <p>
                  <span className="font-semibold">Slug:</span> {project.slug}
                </p>
                <p>
                  <span className="font-semibold">Repository:</span>{" "}
                  {project.repositoryUrl}
                </p>
                {!!project.customDomain && (
                  <p>
                    <span className="font-semibold">Custom Domain:</span>{" "}
                    {project.customDomain}
                  </p>
                )}
              </div>
              <div>
                <p>
                  <span className="font-semibold">Build Command:</span>{" "}
                  {project.buildCommand}
                </p>
                <p>
                  <span className="font-semibold">Created At:</span>{" "}
                  {new Date(project.createdAt).toLocaleString()}
                </p>
                <p>
                  <span className="font-semibold">Domain:</span>{" "}
                  {project.domain}
                </p>
              </div>
            </div>
            <div className="mt-4 flex space-x-2">
              <Tooltip content="Visit the deployed site">
                <Button
                  color="primary"
                  as="a"
                  href={`https://${project.domain}`}
                  target="_blank"
                >
                  Visit Site
                </Button>
              </Tooltip>
              <Tooltip content="Trigger a new build">
                <Button color="secondary">Rebuild</Button>
              </Tooltip>
            </div>
          </CardBody>
        </Card>

        <Card>
          <CardHeader>
            <h2 className="text-xl font-bold">Build Logs</h2>
          </CardHeader>
          <Divider />
          <CardBody>
            <div className="bg-gray-100 p-4 rounded-lg max-h-96 overflow-y-auto">
              {buildLogs.map((log, index) => (
                <div key={index} className="mb-2">
                  <span className="text-gray-500 mr-2">
                    {new Date(log.timestamp).toLocaleTimeString()}
                  </span>
                  <span>{log.message}</span>
                </div>
              ))}
            </div>
          </CardBody>
        </Card>
      </motion.div>
    </div>
  );
};

export default ProjectDetailPage;
