#!/bin/bash
cd /Users/nvictor/go/src/scheduler0-go-client

echo "Committing AccountID type fix..."
git add .
git commit -m "Fix AccountID field type from string to int64 to match private model

- Changed AccountID from string to int64 in Credential struct
- This matches the private model which uses uint64 for accountId
- Fixes JSON unmarshaling error when fetching credentials"

echo "Creating new version tag..."
git tag v1.1.1
git push origin main
git push origin v1.1.1

echo "Published v1.1.1 successfully!"
