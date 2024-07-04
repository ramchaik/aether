import { SQSEvent, SQSRecord, Context } from "aws-lambda";
import { S3 } from "aws-sdk";
import * as childProcess from "child_process";
import * as fs from "fs";
import * as path from "path";
import * as os from "os";

const s3 = new S3();

interface SQSMessageBody {
  repoUrl: string;
  projectId: string;
}

async function runInDocker(repoUrl: string, tempDir: string): Promise<void> {
  // Clone the repository
  await childProcess.execSync(`git clone ${repoUrl} ${tempDir}`);

  // Build and run the Docker container
  await childProcess.execSync(`
    docker build -t secure-build-env .
    docker run --rm \
      --network none \
      --cpus 1 \
      --memory 1g \
      --read-only \
      --tmpfs /tmp \
      -v ${tempDir}:/app \
      secure-build-env
  `);
}

async function deployToS3(
  files: { path: string; content: Buffer }[],
  projectId: string
): Promise<void> {
  console.log(`Deploying ${files.length} files to S3 for project ${projectId}`);

  const bucketName = process.env.S3_BUCKET_NAME!;
  const projectPrefix = `build/${projectId}/`;

  for (const file of files) {
    const params = {
      Bucket: bucketName,
      Key: projectPrefix + path.relative("/tmp/build-", file.path),
      Body: file.content,
    };
    await s3.putObject(params).promise();
  }
}

async function readDirRecursive(
  dir: string
): Promise<{ path: string; content: Buffer }[]> {
  const entries = await fs.promises.readdir(dir, { withFileTypes: true });
  const files = await Promise.all(
    entries.map((entry) => {
      const res = path.resolve(dir, entry.name);
      return entry.isDirectory()
        ? readDirRecursive(res)
        : { path: res, content: fs.readFileSync(res) };
    })
  );
  return files.flat();
}

export async function processMessage(event: SQSEvent, context: Context) {
  for (const record of event.Records) {
    await processRecord(record);
  }
  return {};
}

async function processRecord(record: SQSRecord) {
  console.log("Processing SQS message:", record.body);

  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), "build-"));

  try {
    const messageBody: SQSMessageBody = JSON.parse(record.body);
    if (
      typeof messageBody.repoUrl !== "string" ||
      !messageBody.repoUrl.startsWith("https://")
    ) {
      throw new Error("Invalid repository URL");
    }
    if (
      typeof messageBody.projectId !== "string" ||
      messageBody.projectId.trim() === ""
    ) {
      throw new Error("Invalid project ID");
    }

    await runInDocker(messageBody.repoUrl, tempDir);

    // Load all generated build files
    const distPath = path.join(tempDir, "dist");
    const files = await readDirRecursive(distPath);

    // Deploy to S3 with project-specific prefix
    await deployToS3(files, messageBody.projectId);
    console.log(
      `Successfully deployed to S3 for project ${messageBody.projectId}`
    );
  } catch (error) {
    console.error("Error processing message:", error);
  } finally {
    // Clean up
    fs.rmdirSync(tempDir, { recursive: true });
  }
}
