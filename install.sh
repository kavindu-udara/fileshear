#!/bin/bash

# Build the binary
echo "Building fileshear..."
go build -o fileshear

# Install to /usr/local/bin for system-wide access
echo "Installing fileshear..."
# delete any existing fileshear binary
if [ -f /usr/local/bin/fileshear ]; then
    echo "Removing existing fileshear binary..."
    sudo rm /usr/local/bin/fileshear
fi
# delete any existing fileshear web assets
if [ -d /usr/local/share/fileshear ]; then
    echo "Removing existing fileshear web assets..."
    sudo rm -rf /usr/local/share/fileshear
fi
# Move the binary to /usr/local/bin
echo "Moving fileshear binary to /usr/local/bin..."
sudo mv fileshear /usr/local/bin/

# Copy web assets
# delete any existing web assets
if [ -d /usr/local/share/fileshear ]; then
    echo "Removing existing fileshear web assets..."
    sudo rm -rf /usr/local/share/fileshear
fi
# create the directory for web assets
echo "Copying web assets..."
sudo mkdir -p /usr/local/share/fileshear
sudo cp -r web/ /usr/local/share/fileshear/

# Set permissions
echo "Setting permissions..."
sudo chmod +x /usr/local/bin/fileshear
sudo chmod -R 755 /usr/local/share/fileshear

echo "Installation complete! You can now run 'fileshear' from anywhere."
