import { NextResponse } from "next/server";

const backendApiHost = process.env.BACKEND_API_HOST || "http://localhost:8000";

export async function GET(req: Request) {
  const authHeader = req.headers.get("Authorization");

  const response = await fetch(`${backendApiHost}/api/v1/project`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      ...(authHeader ? { Authorization: authHeader } : {}),
    },
  });
  const data = await response.json();
  return new NextResponse(JSON.stringify(data), { status: 200 });
}

export async function POST(req: Request) {
  const authHeader = req.headers.get("Authorization");
  const response = await fetch(`${backendApiHost}/api/v1/project`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...(authHeader ? { Authorization: authHeader } : {}),
    },
    body: req.body,
  });
  const data = await response.json();
  return new NextResponse(JSON.stringify(data), { status: 200 });
}
