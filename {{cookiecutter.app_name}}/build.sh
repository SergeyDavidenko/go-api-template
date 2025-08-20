#!/bin/bash

# Exit on any error
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check the status of the previous execution
check_result() {
    if [[ $? -ne 0 ]]; then
        print_error "BUILD FAILED"
        exit 1
    fi
}

# Check prerequisites
print_status "Checking prerequisites..."

if ! command_exists docker; then
    print_error "Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command_exists git; then
    print_error "Git is not installed. Please install Git first."
    exit 1
fi

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_warning "Not in a git repository. Version information may be limited."
    HASH="unknown"
    SAVER_TAG="$(date +"%Y%m%d-%H%M")-unknown"
else
    # Get git information
    HASH=$(git rev-parse origin/master 2>/dev/null || git rev-parse HEAD)
    SAVER_TAG="$(date +"%Y%m%d-%H%M")-${HASH:0:7}"
fi

print_status "Build tag: $SAVER_TAG"

# Pull the base image
print_status "Pulling base image..."
docker pull golang:{{cookiecutter.docker_build_image_version}}-alpine
check_result

# Build the Docker image
print_status "Building Docker image..."
docker build -t {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:$SAVER_TAG .
check_result

# Tag as latest
docker tag {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:$SAVER_TAG {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:latest
check_result

print_status "Docker image built successfully: {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:$SAVER_TAG"

# Push to registry
print_status "Pushing to Docker registry..."
docker push {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:$SAVER_TAG
check_result

docker push {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:latest
check_result

print_status "Successfully pushed {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:$SAVER_TAG to registry"

# Print summary
echo ""
print_status "Build Summary:"
echo "  Image: {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:$SAVER_TAG"
echo "  Latest: {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:latest"
echo "  Commit: $HASH"
echo "  Build Time: $(date)"
echo ""

print_status "Build completed successfully! 🎉"