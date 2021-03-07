#!/bin/bash

echo "Installing GO"

sudo apt update && sudo apt upgrade -y
sudo wget https://dl.google.com/go/go1.16.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.16.linux-amd64.tar.gz


