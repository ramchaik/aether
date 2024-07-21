"use client";

import React, { useEffect, useState } from "react";
import {
  Card,
  CardBody,
  CardHeader,
  Divider,
  Chip,
  Button,
  Tooltip,
  Progress,
  CircularProgress,
  ScrollShadow,
} from "@nextui-org/react";
import { AnimatePresence, motion } from "framer-motion";
import { useParams } from "next/navigation";
import {
  deployProject,
  useDeployProject,
  useFetchProject,
  useFetchProjectLogs,
} from "@/hooks/useProjectApi";
import { Project } from "@/store/useProjectStore";
import { CheckIcon, ClockIcon } from "@heroicons/react/24/solid";
import toast from "react-hot-toast";
import Loader from "@/components/loader";
import Error from "@/components/error";

const ProjectDetailPage: React.FC = () => {
  const params = useParams();
  const projectId = params.id as string;
  const {
    data: project,
    isLoading,
    error,
  } = useFetchProject<Project>(projectId);

  const [isPolling, setIsPolling] = useState(false);

  const {
    data: logs,
    isLoading: isLogsLoading,
    error: logsError,
  } = useFetchProjectLogs<any>(projectId, {
    refetchInterval: isPolling ? 5000 : false, // Poll every 5 seconds if isPolling is true
  });

  useEffect(() => {
    if (project?.status === "DEPLOYING") {
      setIsPolling(true);
    } else {
      setIsPolling(false);
    }
  }, [project?.status]);

  const deployProjectMutation = useDeployProject();
  const [isDeploying, setIsDeploying] = useState(false);

  const StatusIndicator = ({ status }: { status: Project["status"] }) => {
    switch (status) {
      case "LIVE":
        return (
          <Chip
            startContent={<CheckIcon className="w-4 h-4" />}
            color="success"
            variant="flat"
          >
            Deployed
          </Chip>
        );
      case "NOT_LIVE":
        return (
          <Chip color="default" variant="flat">
            Not Deployed
          </Chip>
        );
      case "DEPLOYING":
        return (
          <Chip
            startContent={
              <div className="pr-1">
                <CircularProgress
                  size="sm"
                  color="primary"
                  aria-label="Deploying..."
                  classNames={{
                    svg: "w-4 h-4",
                  }}
                />
              </div>
            }
            color="primary"
            variant="flat"
          >
            Deploying
          </Chip>
        );
      default:
        return null;
    }
  };

  const handleRebuild = async () => {
    setIsDeploying(true);
    try {
      await deployProjectMutation.mutateAsync(projectId);
      toast.success("Deployment started successfully!");
    } catch (error) {
      toast.error("Failed to start deployment. Please try again.");
    } finally {
      setIsDeploying(false);
    }
  };

  if (error) return <Error message={"Something went wrong"} />;
  if (isLoading) return <Loader size="lg" color="primary" />;
  if (!project) return <div>Project not found</div>;
  return (
    <div className="container mx-auto px-4 py-10">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <Card className="mb-6">
          <CardHeader className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
            <h1 className="text-2xl font-bold">{project.name}</h1>
            <StatusIndicator status={project.status} />
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
              <Tooltip
                content={
                  project.status === "LIVE"
                    ? "Visit the deployed site"
                    : "Site not yet deployed"
                }
              >
                <Button
                  color="primary"
                  as="a"
                  href={project.domain}
                  target="_blank"
                  isDisabled={project.status !== "LIVE"}
                >
                  Visit Site
                </Button>
              </Tooltip>
              <Tooltip
                content={
                  isDeploying
                    ? "Deployment in progress"
                    : project.status === "LIVE"
                    ? "Trigger a new build"
                    : "Start initial build"
                }
              >
                <Button
                  color="secondary"
                  onClick={handleRebuild}
                  isLoading={isDeploying}
                  isDisabled={isDeploying || project.status === "DEPLOYING"}
                >
                  {isDeploying || project.status === "DEPLOYING"
                    ? "Deploying..."
                    : project.status === "LIVE"
                    ? "Rebuild"
                    : "Build"}
                </Button>
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
            {isLogsLoading ? (
              <Loader size="md" color="primary" />
            ) : logsError ? (
              <Error message="Failed to load logs" />
            ) : logs && logs.length > 0 ? (
              <ScrollShadow className="h-[400px]">
                <div className="bg-gray-100 p-4 rounded-lg">
                  <AnimatePresence>
                    {logs.map((log: any, index: number) => (
                      <motion.div
                        key={index}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -20 }}
                        transition={{ duration: 0.3 }}
                        className="mb-2"
                      >
                        <span className="text-gray-500 mr-2">
                          {new Date(log.timestamp).toLocaleTimeString()}
                        </span>
                        <span>
                          {(() => {
                            try {
                              const parsedLog = JSON.parse(log.log);
                              return parsedLog.stream || log.log;
                            } catch {
                              return log.log;
                            }
                          })()}
                        </span>
                      </motion.div>
                    ))}
                  </AnimatePresence>
                </div>
              </ScrollShadow>
            ) : (
              <p>No logs available</p>
            )}
            {project.status === "DEPLOYING" && (
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{
                  duration: 0.5,
                  repeat: Infinity,
                  repeatType: "reverse",
                }}
                className="mt-4 flex justify-center items-center"
              >
                <CircularProgress size="sm" color="primary" />
                <span className="ml-2">Building in progress...</span>
              </motion.div>
            )}
          </CardBody>
        </Card>
      </motion.div>
    </div>
  );
};

export default ProjectDetailPage;
