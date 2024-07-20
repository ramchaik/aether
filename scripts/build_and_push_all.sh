#!/bin/bash

DIRECTORIES=(
    "./src/forge"
    "./src/frontstage"
    "./src/launchpad"
    "./src/logify"
    "./src/proxy"
)

# aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 502413910473.dkr.ecr.us-east-1.amazonaws.com

for DIR in "${DIRECTORIES[@]}"; do
    echo "Processing directory: $DIR"

    DIR_NAME=$(basename "$DIR")

    cd "$DIR" || {
        echo "Failed to change directory to $DIR"
        continue
    }

    # Run the Docker command with the directory name in the tag

    # docker buildx build --platform linux/amd64,linux/arm64 -t "502413910473.dkr.ecr.us-east-1.amazonaws.com/aether:aether-$DIR_NAME" --push .
    docker buildx build --platform linux/amd64,linux/arm64 -t "docker.io/vsramchaik/aether-$DIR_NAME" --push .

    if [ $? -ne 0 ]; then
        echo "Docker build and push failed in directory: $DIR"
    else
        echo "Docker build and push completed successfully in directory: $DIR"
    fi

    cd - > /dev/null
done
