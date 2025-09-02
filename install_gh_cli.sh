#!/bin/bash

# Install GitHub CLI (gh) - Fixed version
echo "Installing GitHub CLI..."

# Check if already installed
if command -v gh &> /dev/null; then
    echo "✅ GitHub CLI is already installed"
    gh --version
    exit 0
fi

# Install GitHub CLI
type -p wget >/dev/null || (sudo apt update && sudo apt install wget -y)
sudo mkdir -p -m 755 /etc/apt/keyrings
wget -qO- https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null
sudo chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt update
sudo apt install gh -y

if command -v gh &> /dev/null; then
    echo "✅ GitHub CLI installed successfully!"
    gh --version
    echo ""
    echo "To authenticate with GitHub, run:"
    echo "gh auth login"
else
    echo "❌ GitHub CLI installation failed"
    exit 1
fi
