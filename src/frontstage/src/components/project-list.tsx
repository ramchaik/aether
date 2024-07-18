"use client";
import React, { useEffect, useState } from "react";
import {
  Card,
  Button,
  CardHeader,
  CardBody,
  CardFooter,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
} from "@nextui-org/react";
import dynamic from "next/dynamic";
import ProjectForm from "./project-form";
import Link from "next/link";
import { useFetchAllProjects } from "@/hooks/useProjectApi";
import useProjectStore, { Project } from "@/store/useProjectStore";

const MotionDiv = dynamic(
  () => import("framer-motion").then((mod) => mod.motion.div),
  { ssr: false }
);

const ProjectList: React.FC = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { data: projects } = useFetchAllProjects<Project[]>();
  const { setProjects } = useProjectStore();
  useEffect(() => {
    if (projects) {
      setProjects(projects);
    }
  }, [projects, setProjects]);

  const displayProjects = useProjectStore((state) => state.projects);

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">Your Projects</h2>
        <Button color="primary" onPress={() => setIsModalOpen(true)}>
          Create Project
        </Button>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {!!displayProjects.length &&
          displayProjects.map((project) => (
            <MotionDiv
              key={project.id}
              whileHover={{ scale: 1.03 }}
              transition={{ type: "spring", stiffness: 300 }}
            >
              <Card isHoverable className="cursor-pointer">
                <CardHeader className="flex flex-col items-start">
                  <h4 className="text-large font-bold">{project.name}</h4>
                  <p className="text-small text-default-500">
                    {project.repositoryUrl}
                  </p>
                </CardHeader>
                <CardBody>
                  {/* You can add more project details here if needed */}
                </CardBody>
                <CardFooter>
                  <Link
                    href={`/projects/${project.id}`}
                    passHref
                    className="w-full"
                  >
                    <Button as="a" color="primary" className="w-full">
                      View Details
                    </Button>
                  </Link>
                </CardFooter>
              </Card>
            </MotionDiv>
          ))}
      </div>
      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        size="2xl"
      >
        <ModalContent>
          <ModalHeader className="flex flex-col gap-1">
            Create New Project
          </ModalHeader>
          <ModalBody>
            <ProjectForm onClose={() => setIsModalOpen(false)} />
          </ModalBody>
        </ModalContent>
      </Modal>
    </div>
  );
};

export default ProjectList;
