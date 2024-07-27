# Aether 🚀

Aether is a minimalist PaaS (Platform as a Service) offering, inspired by Vercel and Netlify. It allows you to deploy your npm projects effortlessly by simply providing a public repository URL.

## Demo 🎥
![DEMO](docs/aether-demo.gif)

## Features ✨

- Easy deployment from public Git repositories
- Automatic builds and deployments
- Real-time log streaming
- Serverless architecture
- Microservices-based design

## Tech Stack 🛠️

- **Frontend**: NextJS 14, TypeScript, NextUI, Clerk (auth), Tailwind CSS, Zustand, React Query, Framer Motion
- **Backend**: Node.js, Go, Fastify, Chi
- **Database**: PostgreSQL
- **ORM**: Drizzle
- **Message Queue**: AWS SQS
- **Streaming**: AWS Kinesis
- **Storage**: AWS S3
- **Container Orchestration**: Kubernetes (AWS EKS)
- **CI/CD**: GitHub Actions, ArgoCD
- **Infrastructure as Code**: Terraform
- **Cloud Provider**: AWS

## Microservices 🌐

1. **Frontstage**: Public-facing NextJS app for the user interface
2. **Launchpad**: Fastify app for project CRUD operations, uses gRPC for communication with other services
3. **Forge**: Go app for building and deploying projects, uses gRPC for communication with other services
4. **Logify**: Go app for log streaming and aggregation, uses gRPC for communication with other services
5. **Proxy**: Go app serving as a reverse proxy for deployed projects

## Architecture Diagram
![Aether architecture](docs/aether-architecture.png)


## Project Structure 📁

```
├── src
│   ├── frontstage
│   ├── launchpad
│   ├── forge
│   ├── logify
│   └── proxy
├── tf
├── k8s
├── scripts
└── protobuf
```

<!-- ## Getting Started 🏁


TODO: add instruction here -->

## Deployment 🚢

Aether uses GitHub Actions for CI/CD, building and pushing Docker images to DockerHub. ArgoCD is used for GitOps, linking the GitHub repository to the EKS cluster.

## Infrastructure 🏗️

All infrastructure is provisioned using Terraform. Configuration files can be found in the `tf` folder.

## Kubernetes 🎡

Kubernetes manifests for deploying to AWS EKS are located in the `k8s` folder.

## Inter-service Communication 🔀

Aether utilizes gRPC with Protocol Buffers for efficient and type-safe communication between microservices. This ensures high performance and maintainability across the system.

## Contributing 🤝

Contributions are welcome! Please feel free to submit a Pull Request.

## License 📄

This project is licensed under the [MIT License](LICENSE).