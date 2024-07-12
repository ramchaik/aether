import ProjectList from "@/components/project-list";
import { auth, currentUser } from "@clerk/nextjs/server";
import React from "react";
import { Card, CardBody, CardHeader, Avatar, Divider } from "@nextui-org/react";
import ProjectForm from "@/components/project-form";

export default async function DashboardPage() {
  const { userId } = auth();
  const user = await currentUser();

  if (!userId || !user) {
    return (
      <div className="flex items-center justify-center h-screen px-4">
        <Card className="w-full max-w-md">
          <CardBody>
            <p className="text-center text-large">
              You are not logged in. Please login!
            </p>
          </CardBody>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-10">
      <Card className="mb-6">
        <CardHeader className="flex gap-3">
          <Avatar
            src={user.imageUrl}
            size="lg"
            name={`${user.firstName} ${user.lastName}`}
          />
          <div className="flex flex-col">
            <p className="text-md">Welcome,</p>
            <p className="text-xl font-bold">
              {user.firstName} {user.lastName}
            </p>
          </div>
        </CardHeader>
        <Divider />
        <CardBody>
          <div className="space-y-2">
            <p>
              <span className="font-semibold">Email:</span>{" "}
              {user.emailAddresses[0].emailAddress}
            </p>
          </div>
        </CardBody>
      </Card>

      <Card>
        <CardBody>
          <ProjectList />
        </CardBody>
      </Card>
    </div>
  );
}
