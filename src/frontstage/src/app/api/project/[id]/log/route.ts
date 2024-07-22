import { NextResponse } from "next/server";

const logsBackendApiHost =
  process.env.LOGS_BACKEND_API_HOST || "http://0.0.0.0:8080";

console.log({ logsBackendApiHost });

export async function GET(
  req: Request,
  { params }: { params: { id: string } }
) {
  const authHeader = req.headers.get("Authorization");
  const projectId = params.id;

  const response = await fetch(
    `${logsBackendApiHost}/api/v1/${projectId}/logs`,
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
