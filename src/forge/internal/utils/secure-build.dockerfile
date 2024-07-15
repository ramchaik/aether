FROM node:latest

RUN apt-get update && apt-get install -y git

WORKDIR /app

ARG REPO_URL
ARG BUILD_COMMAND

# Clone the repository
RUN git clone ${REPO_URL} ./repo
WORKDIR /app/repo

RUN if [ -f .npmrc ] && grep -q "node-version" .npmrc; then \
    NODE_VERSION=$(grep "node-version" .npmrc | cut -d'=' -f2 | tr -d ' ') && \
    n ${NODE_VERSION} && \
    npm install -g npm@latest; \
    fi

RUN npm i

RUN ${BUILD_COMMAND}

RUN mkdir -p /build

# Move build files to /build directory
RUN mv dist/* /build || mv build/* /build || true
