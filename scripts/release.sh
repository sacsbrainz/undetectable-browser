#!/bin/bash

VERSION=$1
MESSAGE=$2

if [ -z "$VERSION" ]; then
    echo "Error: Version number is required"
    echo "Usage: ./release.sh <version> [message]"
    echo "Example: ./release.sh v1.0.0 'Initial release'"
    exit 1
fi

# Validate version format
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9\.]+)?$ ]]; then
    echo "Error: Version must be in format v1.0.0 or v1.0.0-beta.1"
    exit 1
fi

# If no message is provided, use version as message
if [ -z "$MESSAGE" ]; then
    MESSAGE="Release $VERSION"
fi

# Create and push tag
echo "Creating tag $VERSION..."
git tag -a "$VERSION" -m "$MESSAGE"

echo "Pushing commits and tags..."
git push && git push --tags

echo "Done! GitHub Actions workflow will start automatically."
echo "You can check the progress at: https://github.com/$(git config --get remote.origin.url | sed 's/.*[\/:]\([^\/:]*\/[^\/]*\)\.git/\1/')/actions"