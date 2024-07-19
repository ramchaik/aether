FROM node:20

# Install git and global dependencies
RUN apt-get update && apt-get install -y git && \
    npm install -g npm@latest yarn --force && \
    npm cache clean --force

RUN npm install -g vite @vue/cli @angular/cli \
    create-react-app create-next-app gatsby-cli svelte @sveltejs/kit \
    react-scripts next && \
    npm cache clean --force

WORKDIR /app

ARG REPO_URL
ARG BUILD_COMMAND

# Clone the repository
RUN git clone ${REPO_URL} ./repo
WORKDIR /app/repo

RUN if [ -f yarn.lock ]; then \
    yarn install --frozen-lockfile; \
    elif [ -f package-lock.json ]; then \
    npm ci; \
    else \
    npm install; \
    fi

RUN if grep -q '"react-scripts"' package.json; then \
    npm install react-scripts; \
    fi

# Run the build command
RUN eval ${BUILD_COMMAND}

# Move build files to /build directory
RUN mkdir -p /build && \
    if [ -d build ]; then \
    mv build/* /build; \
    elif [ -d dist ]; then \
    mv dist/* /build; \
    fi