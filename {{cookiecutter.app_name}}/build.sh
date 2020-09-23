#!/bin/bash

HASH=`git rev-parse origin/master`
SAVER_TAG=`date +"%Y%m%d-%H%M"`-${HASH:0:7}

docker pull golang:{{cookiecutter.docker_build_image_version}}
echo "Start build"
docker build -t {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:$SAVER_TAG .
echo "Build docker container {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:$SAVER_TAG"
docker push {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:$SAVER_TAG
echo "Build push container {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:$SAVER_TAG"