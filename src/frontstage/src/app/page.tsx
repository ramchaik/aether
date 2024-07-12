import { Button } from "@nextui-org/react";

export default function Home() {
  return (
    <div className="flex items-center justify-center flex-col gap-10">
      <h1 className="text-4xl font-bold mt-20">aether</h1>
      <h1 className="text-4xl">Welcome to aether</h1>
      <div>
        <Button>Click me</Button>
      </div>
    </div>
  );
}
