#!/bin/bash

function check_docker_installed() {
    if ! command -v docker &>/dev/null; then
        return 1
    else
        return 0
    fi
}

# Check if Docker is already installed
if check_docker_installed; then
    echo "Docker is already installed on vm."
else
    echo "Docker is not installed. Installing Docker..."

    # Update the Package Repository
    sudo apt update -y
    # Install Prerequisite Packages
    sudo apt install apt-transport-https ca-certificates curl software-properties-common -y

    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    # Add Docker Repository
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    # Specify Installation Source
    apt-cache policy docker-ce
    # Install Docker
    sudo apt install docker-ce -y

    echo "Docker has been successfully installed."
fi

#sudo systemctl status docker
