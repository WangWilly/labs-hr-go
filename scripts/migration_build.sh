#/bin/bash

IMG_NAME="labs-hr-go-migration:latest"

# Build docker image
docker build -t $IMG_NAME -f database/Dockerfile .
if [ $? -ne 0 ]; then
    echo "Failed to build the Docker image. Please check the errors above."
    exit 1
fi

echo "Docker image built successfully. Image name: $IMG_NAME:latest"
