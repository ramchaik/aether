# Aether 🚀

Aether is a minimalist PaaS (Platform as a Service) offering, inspired by Vercel and Netlify. It allows you to deploy your npm projects effortlessly by simply providing a public repository URL.

## Architecture Diagram
![Aether architecture](docs/aether-architecture.png)

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
2. **Launchpad**: Fastify app for project CRUD operations
3. **Forge**: Go app for building and deploying projects
4. **Logify**: Go app for log streaming and aggregation
5. **Proxy**: Go app serving as a reverse proxy for deployed projects

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
└── scripts
```

<!-- ## Demo 🎥
TODO: add demo  -->

<!-- ## Getting Started 🏁

TODO: add instruction here -->

## Deployment 🚢

Aether uses GitHub Actions for CI/CD, building and pushing Docker images to DockerHub. ArgoCD is used for GitOps, linking the GitHub repository to the EKS cluster.

## Infrastructure 🏗️

All infrastructure is provisioned using Terraform. Configuration files can be found in the `tf` folder.

## Kubernetes 🎡

Kubernetes manifests for deploying to AWS EKS are located in the `k8s` folder.

## Contributing 🤝

Contributions are welcome! Please feel free to submit a Pull Request.

## License 📄

This project is licensed under the [MIT License](LICENSE).
