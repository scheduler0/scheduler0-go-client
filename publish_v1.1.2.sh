#!/bin/bash

# Publish scheduler0-go-client v1.1.2
echo "Publishing scheduler0-go-client v1.1.2..."

# Add all changes
git add .

# Commit changes
git commit -m "Fix Project ID types from string to int64 to match backend API

- Updated Project struct to use int64 for ID and AccountID fields
- Updated GetProject, UpdateProject, DeleteProject methods to accept int64
- Fixed JSON unmarshaling issues with backend API responses
- Backend returns IDs as numbers, not strings"

# Create and push tag
git tag v1.1.2
git push origin main
git push origin v1.1.2

echo "Published v1.1.2 successfully!"
