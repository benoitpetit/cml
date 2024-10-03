#!/bin/bash

# Script to uninstall CML from the system
set -e

# Define the binary name
BINARY_NAME="cml"

# Function to remove CML on Debian-based systems
uninstall_debian_based() {
    echo "Uninstalling CML from Debian-based system..."

    # Check if the binary exists
    if [ -f /usr/local/bin/$BINARY_NAME ]; then
        sudo rm /usr/local/bin/$BINARY_NAME
        echo "CML has been successfully uninstalled."
    else
        echo "CML is not installed."
    fi
}

# Function to remove CML on Red Hat-based systems
uninstall_redhat_based() {
    echo "Uninstalling CML from Red Hat-based system..."

    # Check if the binary exists
    if [ -f /usr/local/bin/$BINARY_NAME ]; then
        sudo rm /usr/local/bin/$BINARY_NAME
        echo "CML has been successfully uninstalled."
    else
        echo "CML is not installed."
    fi
}

# Function to remove CML on openSUSE
uninstall_opensuse() {
    echo "Uninstalling CML from openSUSE..."

    # Check if the binary exists
    if [ -f /usr/local/bin/$BINARY_NAME ]; then
        sudo rm /usr/local/bin/$BINARY_NAME
        echo "CML has been successfully uninstalled."
    else
        echo "CML is not installed."
    fi
}

# Check the OS and uninstall accordingly
if [ -f /etc/os-release ]; then
    . /etc/os-release
    case "$ID_LIKE" in
        debian)
            uninstall_debian_based
            ;;
        rhel|fedora)
            uninstall_redhat_based
            ;;
        suse*)
            uninstall_opensuse
            ;;
        *)
            echo "Unsupported Linux distribution: $ID"
            exit 1
            ;;
    esac
else
    echo "Unable to detect Linux distribution."
    exit 1
fi
