// ProjectForm.tsx

"use client";

import React, { useState } from "react";
import { Input, Button, Textarea, Tooltip } from "@nextui-org/react";
import { motion } from "framer-motion";

interface ProjectFormProps {
  onClose: () => void;
}

const ProjectForm: React.FC<ProjectFormProps> = ({ onClose }) => {
  const [formData, setFormData] = useState({
    projectName: "",
    repoUrl: "",
    customDomain: "",
    buildCommand: "",
  });

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log("Form submitted:", formData);
    // Here you would typically send the data to your backend
    onClose(); // Close the modal after submission
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Project Name"
          name="projectName"
          value={formData.projectName}
          onChange={handleInputChange}
          required
        />
        <Input
          label="Repository URL"
          name="repoUrl"
          value={formData.repoUrl}
          onChange={handleInputChange}
          required
        />
        <Input
          label="Custom Domain (optional)"
          name="customDomain"
          value={formData.customDomain}
          onChange={handleInputChange}
        />
        <Tooltip content="e.g., npm run dev, npm start">
          <Input
            label="Build Command"
            name="buildCommand"
            value={formData.buildCommand}
            onChange={handleInputChange}
            required
          />
        </Tooltip>
        <div className="flex justify-end gap-2 mt-4">
          <Button color="danger" variant="light" onPress={onClose}>
            Cancel
          </Button>
          <motion.div whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
            <Button type="submit" color="primary">
              Create Project
            </Button>
          </motion.div>
        </div>
      </form>
    </motion.div>
  );
};

export default ProjectForm;
