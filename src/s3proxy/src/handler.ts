import { APIGatewayProxyHandler } from "aws-lambda";

export const proxy: APIGatewayProxyHandler = async (event) => {
  return {
    statusCode: 200,
    body: JSON.stringify({
      message: "Hello, from Proxy!",
      input: event,
    }),
  };
};
