#!/bin/bash

WORK_DIR=`mktemp -d`

echo ${WORK_DIR}

cd ${WORK_DIR} || exit 1

function usuage () {
  echo "Usage: script.sh -p <file-output-path> -e <envs> -start-cmd <start-cmd> -build-cmd <build-cmd> -app-name <app-name>"
}

while getopts "e:b:s:n:l:" opt
do
  case ${opt} in
    s)
       start_cmd=${OPTARG}
       ;;
    b)
       build_cmd=${OPTARG}
       ;;
    n)
      app_name=${OPTARG}
      ;;
    e)
      envs=${OPTARG}
      ;;
    l)
      repo_url=${OPTARG}
      ;;
    ?)
       echo "Invalid option: -${OPTARG}."
       exit 1
       ;;
  esac
done


# Check if required arguments are provided
if [[  -z $envs || -z $start_cmd || -z $build_cmd || -z $app_name ]]
then
    usuage
    exit 1
fi

if [[ ! -z $repo_url ]]
then
  echo "🚕 Pulling repository..."
  git clone "${repo_url}"

  if [[ $? -eq 0 ]]
  then
    echo "✅ Repository Pulled Successfully"
  else
    echo "❌ Failed To Pull Repository Code. please confirm wave-deploy has access to the repository on github"
    exit 1
  fi
fi

# Build the application
echo "👷🏽‍Building application"
nixpacks build . \
  --start-cmd "${start_cmd}" \
  --build-cmd "${build_cmd}" \
  --platform linux/amd64 \
  --name "${app_name}" \
  --out "${WORK_DIR}" &> /dev/null

if [[ $? -eq 0 ]]
then
  echo "✅ Application Build Successful"
else
  echo "❌ Application Build Failed ❌"
fi

# TODO:: On install of cli, make sure all used dependencies are available on host (git, buildpack)
# TODO:: Remember to delete directory once all operations the directory is needed for has been completed
# TODO:: Add spinner
echo "$WORK_DIR"