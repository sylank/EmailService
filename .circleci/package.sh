#!/bin/bash

CONFIG_DIR="./config"

VERSION=$(cat infrastructure/VERSION)
APP_NAME=$(cat infrastructure/APP_NAME)

echo "Download credentials..."
rm $CONFIG_DIR/secrets.json

aws s3 cp s3://artifactory.levendulabalatonmaria.info/credentials/$APP_NAME/secrets.json $CONFIG_DIR/secrets.json

ls -la $CONFIG_DIR