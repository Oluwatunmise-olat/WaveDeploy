#!/bin/bash

function check_docker_installed() {
    if ! command -v docker &>/dev/null; then
        return 1
    else
        return 0
    fi
}

function check_nixpacks_installed() {
    if ! command -v nixpacks &>/dev/null; then
        return 1
    else
        return 0
    fi
}

function check_caddy_installed() {
    if ! command -v caddy &>/dev/null; then
        return 1
    else
        return 0
    fi
}

if check_docker_installed; then
    echo "Docker is already installed on vm."
else
    echo "Installing Docker..."

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


if check_nixpacks_installed; then
    echo "Nixpacks is already installed."
else
    echo "Installing Nixpacks..."

    # Install Nixpacks
    curl -sSL https://nixpacks.com/install.sh | bash

    if check_nixpacks_installed; then
        echo "Nixpacks has been successfully installed."
    else
        echo "Failed to install Nixpacks."
    fi
fi


if check_caddy_installed; then
    echo "Caddy is already installed."
else
    echo "Installing Caddy..."

    sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https curl
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
    sudo apt update
    sudo apt install caddy

    if check_caddy_installed; then
        echo "Caddy has been successfully installed."
    else
        echo "Failed to install Caddy."
    fi
fi
