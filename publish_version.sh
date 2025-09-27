#!/bin/bash
cd /Users/nvictor/go/src/scheduler0-go-client

# Check current version
echo "Current tags:"
git tag --list | grep v1 | tail -5

# Get the latest version number
LATEST_VERSION=$(git tag --list | grep v1 | sort -V | tail -1)
echo "Latest version: $LATEST_VERSION"

# Extract version number and increment
if [[ $LATEST_VERSION =~ v1\.([0-9]+) ]]; then
    VERSION_NUM=${BASH_REMATCH[1]}
    NEW_VERSION_NUM=$((VERSION_NUM + 1))
    NEW_VERSION="v1.$NEW_VERSION_NUM"
else
    NEW_VERSION="v1.0.3"
fi

echo "Creating new version: $NEW_VERSION"

# Add and commit changes
git add .
git commit -m "Update Credential struct with camelCase JSON tags and additional fields

- Changed api_key to apiKey
- Changed api_secret to apiSecret  
- Changed date_created to dateCreated
- Added missing fields: accountId, dateModified, dateDeleted, createdBy, modifiedBy, deletedBy
- Updated field types to match private models and OpenAPI schema"

# Create and push tag
git tag $NEW_VERSION
git push origin main
git push origin $NEW_VERSION

echo "Published $NEW_VERSION successfully!"
echo "Update app.scheduler0.com to use: github.com/scheduler0/scheduler0-go-client $NEW_VERSION"
