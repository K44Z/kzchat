
#!/bin/bash

set -e

APP_NAME="kzchat"
SERVER_NAME="kzchat-server"
CLIENT_BINARY_PATH="./cmd/tui"
SERVER_BINARY_PATH="./cmd/server"
CONFIG_DIR="/etc/$APP_NAME"
CONFIG_DEST="$CONFIG_DIR/.env"
CLIENT_INSTALL_PATH="/usr/local/bin/$APP_NAME"
SERVER_INSTALL_PATH="/usr/local/bin/$SERVER_NAME"

read -p "Database URL [postgres://user:password@localhost:5432/dbname?sslmode=disable]: " DB_URL
DB_URL=${DB_URL:-"postgres://user:password@localhost:5432/dbname?sslmode=disable"}

read -p "Application port [:8080]: " PORT
PORT=${PORT:-":8080"}

JWT_SECRET=$(openssl rand -hex 32)

go build -o $APP_NAME $CLIENT_BINARY_PATH
sudo mv $APP_NAME $CLIENT_INSTALL_PATH
sudo chmod +x $CLIENT_INSTALL_PATH

go build -o $SERVER_NAME $SERVER_BINARY_PATH
sudo mv $SERVER_NAME $SERVER_INSTALL_PATH
sudo chmod +x $SERVER_INSTALL_PATH

sudo mkdir -p $CONFIG_DIR

sudo tee "$CONFIG_DEST" > /dev/null <<EOF
DB_URL=$DB_URL
PORT=$PORT
JWT_SECRET=$JWT_SECRET
EOF

echo "Installation complete"
echo "Config saved at: $CONFIG_DEST"
echo "Run client with: $APP_NAME"
echo "Run server with: $SERVER_NAME"

