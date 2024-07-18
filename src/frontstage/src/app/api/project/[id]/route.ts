import { NextResponse } from "next/server";

const backendApiHost = process.env.BACKEND_API_HOST || "http://0.0.0.0:8000";

console.log({ backendApiHost });

export async function GET(
  req: Request,
  { params }: { params: { id: string } }
) {
  const authHeader = req.headers.get("Authorization");
  const projectId = params.id;

  const response = await fetch(
    `${backendApiHost}/api/v1/project/${projectId}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        ...(authHeader ? { Authorization: authHeader } : {}),
      },
    }
  );
  const data = await response.json();
  return new NextResponse(JSON.stringify(data), { status: 200 });
}
