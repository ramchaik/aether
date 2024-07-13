"use client";
import { Card, Button, Image } from "@nextui-org/react";
import { motion } from "framer-motion";
import Link from "next/link";

export default function Home() {
  return (
    <Card className="w-full h-screen flex flex-col items-center justify-center bg-gradient-to-br from-blue-900 to-purple-900">
      <motion.div
        initial={{ opacity: 0, y: -50 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 1 }}
        className="text-center"
      >
        <h1 className="text-6xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-600">
          aether
        </h1>
      </motion.div>

      <motion.div
        initial={{ opacity: 0, scale: 0.5 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 1, delay: 0.5 }}
        className="my-8"
      >
        <Image src="rocket.svg" alt="Rocket" width={96} height={96} />
      </motion.div>

      <motion.div
        initial={{ opacity: 0, y: 50 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 1, delay: 1 }}
        className="text-center"
      >
        <h2 className="text-2xl text-white mb-4">Welcome to aether</h2>
        <p className="text-lg text-white max-w-md mx-auto mb-8">
          A simple way to deploy your frontend applications with fast delivery.
        </p>
      </motion.div>

      <Link href="/sign-in">
        <Button color="primary" variant="shadow" size="lg">
          Get Started
        </Button>
      </Link>
    </Card>
  );
}
