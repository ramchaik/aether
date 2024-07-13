FROM node:20
ARG DATABASE_URL
ENV DATABASE_URL=${DATABASE_URL}

WORKDIR /app

RUN npm i -g pnpm

COPY package*.json ./
RUN pnpm install

# Copy all source files
COPY . .

# Build the application
RUN pnpm run build

EXPOSE 8000
CMD [ "pnpm", "start" ]
