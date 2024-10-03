#!/bin/bash

# Function to display error messages and exit
function error_exit {
    echo "$1" 1>&2
    exit 1
}

# Create release directories
mkdir -p releases/linux/amd64
mkdir -p releases/windows/amd64

# Compile for Linux
echo "Compiling for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o releases/linux/amd64/cml main.go
if [ $? -ne 0 ]; then
  error_exit "Error during compilation for Linux."
fi

# Add execute permissions for the Linux executable
chmod +x releases/linux/amd64/cml
if [ $? -ne 0 ]; then
  error_exit "Error adding execute permissions to the Linux executable."
fi

# Compress the Linux executable into tar.gz
echo "Compressing the Linux executable into cml-linux-amd64.tar.gz..."
tar -czvf releases/linux/amd64/cml-linux-amd64.tar.gz -C releases/linux/amd64 cml
if [ $? -ne 0 ]; then
  error_exit "Error compressing the Linux executable."
fi

# Compile for Windows
echo "Compiling for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o releases/windows/amd64/cml.exe main.go
if [ $? -ne 0 ]; then
  error_exit "Error during compilation for Windows."
fi

# Check if zip is installed
if ! command -v zip &> /dev/null
then
    error_exit "The 'zip' program is not installed. Please install it to continue."
fi

# Compress the Windows executable into zip
echo "Compressing the Windows executable into cml-windows-amd64.zip..."
zip -j releases/windows/amd64/cml-windows-amd64.zip releases/windows/amd64/cml.exe
if [ $? -ne 0 ]; then
  error_exit "Error compressing the Windows executable."
fi

echo "Compilation and compression completed successfully!"
