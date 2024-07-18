import {
  SQSClient,
  SendMessageCommand,
  MessageAttributeValue,
} from "@aws-sdk/client-sqs";

const awsRegion = process.env.AWS_REGION;
const queueUrl = process.env.AWS_QUEUE_URL;

// Configure the SQS client
const sqsClient = new SQSClient({ region: awsRegion });

interface MessageAttributes {
  [key: string]: MessageAttributeValue;
}

export async function pushMessageToDeployQueue(
  message: unknown,
  messageType: string = "Build"
) {
  if (!queueUrl) {
    throw new Error(
      "AWS_QUEUE_URL is not defined in the environment variables"
    );
  }

  const messageAttributes: MessageAttributes = {
    MessageType: {
      DataType: "String",
      StringValue: messageType,
    },
  };

  try {
    const command = new SendMessageCommand({
      QueueUrl: queueUrl,
      MessageBody: JSON.stringify(message),
      MessageAttributes: messageAttributes,
    });

    const response = await sqsClient.send(command);

    if (!response.MessageId) {
      throw new Error("Failed to get MessageId from SQS response");
    }

    return response.MessageId;
  } catch (error) {
    console.error("Error sending message to SQS:", error);
    throw error; // Re-throw the error to allow the caller to handle it
  }
}
