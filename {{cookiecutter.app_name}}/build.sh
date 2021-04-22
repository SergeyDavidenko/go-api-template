#!/bin/bash


# To check the status of the previous execution
function check_result {
  if [[ "$?" -ne 0 ]] ; then
          echo "BUILD FAILED"
          exit 1
  fi
}


HASH=`git rev-parse origin/master`
SAVER_TAG=`date +"%Y%m%d-%H%M"`-${HASH:0:7}

docker pull golang:{{cookiecutter.docker_build_image_version}}
check_result
echo "Start build"
docker build -t {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:$SAVER_TAG .
check_result
echo "Build docker container {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:$SAVER_TAG"
docker push {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:$SAVER_TAG
check_result
echo "Build push container {{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}:$SAVER_TAG"