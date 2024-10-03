#!/bin/bash

# Script to install CML on various Linux distributions
set -e

# Define the binary name and the latest release URL
BINARY_NAME="cml"
VERSION="v1.0.0"
RELEASE_URL="https://github.com/benoitpetit/cml/releases/download/${VERSION}" # Updated to use the correct version URL

# Function to install CML on Debian-based systems
install_debian_based() {
    echo "Installing CML on a Debian-based system..."
    sudo apt update
    sudo apt install -y wget

    # Download the binary
    echo "Downloading CML..."
    wget -O /tmp/$BINARY_NAME.tar.gz $RELEASE_URL/$BINARY_NAME-linux-amd64.tar.gz

    # Extract and move the binary
    echo "Installing CML..."
    sudo tar -xzvf /tmp/$BINARY_NAME.tar.gz -C /usr/local/bin/
    sudo chmod +x /usr/local/bin/$BINARY_NAME

    # Clean up
    rm /tmp/$BINARY_NAME.tar.gz
    echo "CML installation completed successfully!"
}

# Function to install CML on Red Hat-based systems
install_redhat_based() {
    echo "Installing CML on a Red Hat-based system..."
    sudo dnf install -y wget || sudo yum install -y wget

    # Download the binary
    echo "Downloading CML..."
    wget -O /tmp/$BINARY_NAME.tar.gz $RELEASE_URL/$BINARY_NAME-linux-amd64.tar.gz

    # Extract and move the binary
    echo "Installing CML..."
    sudo tar -xzvf /tmp/$BINARY_NAME.tar.gz -C /usr/local/bin/
    sudo chmod +x /usr/local/bin/$BINARY_NAME

    # Clean up
    rm /tmp/$BINARY_NAME.tar.gz
    echo "CML installation completed successfully!"
}

# Function to install CML on openSUSE
install_opensuse() {
    echo "Installing CML on openSUSE..."
    sudo zypper install -y wget

    # Download the binary
    echo "Downloading CML..."
    wget -O /tmp/$BINARY_NAME.tar.gz $RELEASE_URL/$BINARY_NAME-linux-amd64.tar.gz

    # Extract and move the binary
    echo "Installing CML..."
    sudo tar -xzvf /tmp/$BINARY_NAME.tar.gz -C /usr/local/bin/
    sudo chmod +x /usr/local/bin/$BINARY_NAME

    # Clean up
    rm /tmp/$BINARY_NAME.tar.gz
    echo "CML installation completed successfully!"
}

# Check the OS and install accordingly
if [ -f /etc/os-release ]; then
    . /etc/os-release
    case "$ID_LIKE" in
        debian)
            install_debian_based
            ;;
        rhel|fedora)
            install_redhat_based
            ;;
        suse*)
            install_opensuse
            ;;
        *)
            echo "Unsupported Linux distribution: $ID"
            echo "Please install CML manually."
            exit 1
            ;;
    esac
else
    echo "Unable to detect Linux distribution."
    exit 1
fi

