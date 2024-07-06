FROM node:20

WORKDIR /app

RUN apt-get update && apt-get install -y git

RUN useradd -m builduser
USER builduser

VOLUME /app

CMD ["sh", "-c", "npm install && npm run build"]