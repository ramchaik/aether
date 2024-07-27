import React from "react";
import { Spinner } from "@nextui-org/react";
import { motion } from "framer-motion";

interface LoaderProps {
  size?: "sm" | "md" | "lg";
  color?: "primary" | "secondary" | "success" | "warning" | "danger";
}

const Loader: React.FC<LoaderProps> = ({ size = "md", color = "primary" }) => {
  return (
    <div className="fixed inset-0 flex items-center justify-center bg-white bg-opacity-75 z-[9999] pt-16">
      <motion.div
        initial={{ opacity: 0, scale: 0.5 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.5 }}
        className="flex flex-col items-center justify-center"
      >
        <Spinner size={size} color={color} />
        <motion.p
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2, duration: 0.5 }}
          className="mt-4 text-center text-gray-600"
        >
          Loading...
        </motion.p>
      </motion.div>
    </div>
  );
};

export default Loader;
