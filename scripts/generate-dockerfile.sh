#!/bin/bash

usage() {
  echo "Usage: wave-deploy [-e <env1=value1> [-e <env2=value2> ...]] -s <start-cmd> -b <build-cmd> -n <app-name> -w <work directory> -p <app path> -o <output path>"
}

# Parse options
envs=()
while getopts ":e:b:p:s:n:l:w:o:" opt; do
  case $opt in
    s) start_cmd=$OPTARG ;;
    b) build_cmd=$OPTARG ;;
    n) app_name=$OPTARG ;;
    e) envs+=("$OPTARG") ;;
    l) repo_url=$OPTARG ;;
    w) work_dir=$OPTARG ;;
    p) app_path=$OPTARG ;;
    o) output_path=$OPTARG ;;
    ?)
      echo "Invalid option: -$OPTARG"
      usage
      exit 1
      ;;
  esac
done

# Check if required arguments are provided
if [[ (-z $app_name) || (-z $work_dir) || (-z $app_path) || (-z $output_path) ]]; then
  usage
  exit 1
fi

# Navigate to the work directory
cd $work_dir || exit 1

# Clone repository if provided
if [[ ! -z $repo_url ]]; then
  echo "🚕 Pulling repository..."
  git clone "$repo_url"

  if [[ $? -eq 0 ]]; then
    echo "✅ Repository Pulled Successfully"
  else
    echo "❌ Failed To Pull Repository Code. Please confirm wave-deploy has access to the repository on GitHub"
    exit 1
  fi
fi

echo "👷🏽 Building application"

env_args=()
for env in "${envs[@]}"; do
  env_args+=("--env" "$env")
done

# Include build and start commands if provided
if [[ ! -z $build_cmd ]]; then
  build_args+=( "--build-cmd" "$build_cmd" )
fi

if [[ ! -z $start_cmd ]]; then
  build_args+=( "--start-cmd" "$start_cmd" )
fi

# Build the application
nixpacks build "${app_path}" \
  --platform linux/amd64 \
  --name "$app_name" \
   "${build_args[@]}" \
  "${env_args[@]}" \
  --out "${output_path}" &> /dev/null

mv "${output_path}/.nixpacks/Dockerfile" "${output_path}/.nixpacks/Dockerfile.wavedeploy"
mv "${output_path}/.nixpacks/Dockerfile.wavedeploy" "${app_path}/"
mv "${output_path}/.nixpacks" "${app_path}/"

if [[ $? -eq 0 ]]; then
  echo "✅ Application Build Successful"
else
  echo "❌ Application Build Failed ❌"
  exit 1
fi

# TODO: On install of cli, make sure all used dependencies are available on host (git, buildpack)
# TODO: Remember to delete directory once all operations the directory is needed for has been completed
# TODO: Add spinner
exit 0
