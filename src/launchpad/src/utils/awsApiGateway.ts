import {
  APIGatewayClient,
  GetRestApisCommand,
  GetStageCommand,
} from "@aws-sdk/client-api-gateway";

const awsRegion = process.env.AWS_REGION;
const serviceName = "s3proxy";
const serviceStage = "dev";

let ApiGatewayUrl: null | string = null;

async function findApiId(
  serviceName: string,
  stage: string
): Promise<string | null> {
  const client = new APIGatewayClient({ region: awsRegion });
  const command = new GetRestApisCommand({});

  try {
    const response = await client.send(command);
    const api = response.items?.find(
      (item) => item.name === `${stage}-${serviceName}`
    );
    return api?.id || null;
  } catch (error) {
    console.error("Error finding API ID:", error);
    return null;
  }
}

export async function fetchApiUrl(): Promise<string> {
  if (!!ApiGatewayUrl) return ApiGatewayUrl;

  const apiId = await findApiId(serviceName, serviceStage);
  if (!apiId) {
    throw new Error("Could not find API ID");
  }

  const client = new APIGatewayClient({ region: awsRegion });
  const command = new GetStageCommand({
    restApiId: apiId,
    stageName: "dev",
  });

  try {
    const response = await client.send(command);
    ApiGatewayUrl = `https://${apiId}.execute-api.${awsRegion}.amazonaws.com/${response.stageName}`;
    return ApiGatewayUrl;
  } catch (error) {
    console.error("Error fetching API Gateway URL:", error);
    throw new Error("Failed to fetch API Gateway URL");
  }
}
