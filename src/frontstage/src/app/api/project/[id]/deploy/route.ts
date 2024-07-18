import { NextResponse } from "next/server";

const backendApiHost = process.env.BACKEND_API_HOST || "http://0.0.0.0:8000";

console.log({ backendApiHost });

export async function POST(
  req: Request,
  { params }: { params: { id: string } }
) {
  console.log("POST request on the next BE");
  const authHeader = req.headers.get("Authorization");
  const requestBody = JSON.stringify(await req.json());
  const projectId = params.id;

  const response = await fetch(
    `${backendApiHost}/api/v1/project/${projectId}/deploy`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...(authHeader ? { Authorization: authHeader } : {}),
      },
      body: requestBody,
    }
  );
  const data = await response.json();
  return new NextResponse(JSON.stringify(data), { status: 200 });
}
