#!/bin/bash

###############################################################################################################################

#################################################### Creating user and group ####################################################

###############################################################################################################################

sudo groupadd csye6225
sudo useradd -s /bin/false -g csye6225 -d /opt/app -m webapp

###############################################################################################################################

#################################################### Moving files from tmp ####################################################

###############################################################################################################################

sudo touch /opt/app/.env
sudo chown admin:admin /opt/app/.env
sudo mv /tmp/users.csv /opt/
sudo mv /tmp/assessment-application /opt/app/
sudo mv /tmp/assessment-application.service /etc/systemd/system/

###############################################################################################################################

######################################################### Golang Setup ########################################################

###############################################################################################################################

# Update package list and install dependencies

sudo apt-get update -y
sudo apt-get upgrade -y

# Define the Go version to install
GO_VERSION="1.21.1"
GO_ARCHIVE="go$GO_VERSION.linux-amd64.tar.gz"

# Download and extract the Go binary distribution
echo "Downloading and installing Go $GO_VERSION..."
wget https://dl.google.com/go/$GO_ARCHIVE
sudo tar -C /usr/local -xzf $GO_ARCHIVE
rm -f $GO_ARCHIVE

# Set Go environment variables
echo "Setting Go environment variables..."
echo "export GOROOT=/usr/local/go" >> ~/.bashrc
echo "export GOPATH=\$HOME/go" >> ~/.bashrc
echo "export PATH=\$PATH:\$GOROOT/bin:\$GOPATH/bin" >> ~/.bashrc

# Load the updated environment variables
source ~/.bashrc

###############################################################################################################################

#################################################### Enabling service in Systemd ##############################################

###############################################################################################################################

sudo systemctl daemon-reload
sudo systemctl enable assessment-application