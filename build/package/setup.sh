#!/bin/bash

###############################################################################################################################

#################################################### Moving files from tmp ####################################################

###############################################################################################################################

sudo mv /tmp/users.csv /opt/
mv /tmp/assessment-application /home/admin/

###############################################################################################################################

######################################################### Golang Setup ########################################################

###############################################################################################################################

# Update package list and install dependencies

sudo apt-get update -y
sudo apt-get upgrade -y

# Define the Go version to install
GO_VERSION="1.21.1"
GO_ARCHIVE="go$GO_VERSION.linux-amd64.tar.gz"

# Define database credentials
username="pranay"
password="password"

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

######################################################## Database Setup #######################################################

###############################################################################################################################

# Update package list and install MariaDB
echo "Updating package list and installing MariaDB..."
sudo apt update
sudo apt install -y mariadb-server

# Start MariaDB service and enable it to start on boot
echo "Starting MariaDB service..."
sudo systemctl start mariadb
sudo systemctl enable mariadb

# Secure the MariaDB installation (set root password and remove test databases)
echo "Securing MariaDB installation..."
sudo mysql_secure_installation <<EOF

y
$password
$password
y
y
y
y
EOF

# Create a new database and the specified user with the provided password
echo "Creating MariaDB user '$username'..."
sudo mysql -u root -p$password <<EOF
CREATE USER '$username'@'localhost' IDENTIFIED BY '$password';
GRANT ALL PRIVILEGES ON *.* TO '$username'@'localhost';
FLUSH PRIVILEGES;
EOF


# Save the password and username in the user's profile
echo "export DB_USER='$username'" >> .env  # For bash
echo "export DB_PASSWORD='$password'" >> .env  # For bash

# Reload the profile
source ~/.bashrc

echo "MariaDB installation, user setup, and username/password saved in profile complete."