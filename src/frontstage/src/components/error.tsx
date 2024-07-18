import React from "react";
import { Button } from "@nextui-org/react";
import { motion } from "framer-motion";

interface ErrorProps {
  message: string;
  onRetry?: () => void;
}

const Error: React.FC<ErrorProps> = ({ message, onRetry }) => {
  return (
    <div className="fixed inset-0 flex items-center justify-center bg-white bg-opacity-75 z-[9999] pt-16">
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="flex flex-col items-center justify-center max-w-md p-8 bg-white rounded-lg shadow-lg"
      >
        <motion.svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth={1.5}
          stroke="currentColor"
          className="w-16 h-16 text-red-500 mb-4"
          initial={{ scale: 0.5, rotate: -180 }}
          animate={{ scale: 1, rotate: 0 }}
          transition={{ duration: 0.5 }}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z"
          />
        </motion.svg>

        <h2 className="text-2xl font-bold text-gray-800 mb-4">Error</h2>
        <p className="text-gray-600 text-center mb-6">{message}</p>
        {onRetry && (
          <Button color="primary" onClick={onRetry} className="font-semibold">
            Try Again
          </Button>
        )}
      </motion.div>
    </div>
  );
};

export default Error;
