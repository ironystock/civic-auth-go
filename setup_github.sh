#!/bin/bash

# Civic Auth Go SDK - GitHub Setup Script

echo "🚀 Setting up GitHub repository for Civic Auth Go SDK"
echo "=================================================="

# Check if gh CLI is available
if command -v gh &> /dev/null; then
    echo "✅ GitHub CLI detected"
    
    read -p "Create a new GitHub repository? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Creating GitHub repository..."
        
        # Create repository with gh CLI
        gh repo create civic-auth-go \
            --description "A comprehensive Go SDK for integrating with Civic Auth's OIDC/OAuth2 authentication service" \
            --public \
            --source=. \
            --remote=origin \
            --push
        
        if [ $? -eq 0 ]; then
            echo "✅ Repository created and pushed successfully!"
            echo "🌐 Repository URL: https://github.com/$(gh api user --jq .login)/civic-auth-go"
        else
            echo "❌ Failed to create repository with gh CLI"
            echo "Please create the repository manually on GitHub"
        fi
    else
        echo "Skipped repository creation"
    fi
else
    echo "GitHub CLI not found. Please follow manual setup instructions below:"
    echo ""
    echo "MANUAL SETUP INSTRUCTIONS:"
    echo "========================="
    echo "1. Go to https://github.com/new"
    echo "2. Repository name: civic-auth-go"
    echo "3. Description: A comprehensive Go SDK for integrating with Civic Auth's OIDC/OAuth2 authentication service"
    echo "4. Make it Public"
    echo "5. Don't initialize with README, .gitignore, or license (we already have them)"
    echo "6. Click 'Create repository'"
    echo "7. Then run these commands:"
    echo ""
    echo "   git remote add origin https://github.com/YOUR_USERNAME/civic-auth-go.git"
    echo "   git push -u origin main"
    echo ""
fi

echo ""
echo "📋 Repository Information:"
echo "Name: civic-auth-go"
echo "Description: A comprehensive Go SDK for integrating with Civic Auth's OIDC/OAuth2 authentication service"
echo "Language: Go"
echo "License: MIT"
echo "Topics: civic, auth, oidc, oauth2, go, sdk, authentication"
