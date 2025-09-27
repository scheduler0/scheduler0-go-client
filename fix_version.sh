#!/bin/bash
cd /Users/nvictor/go/src/scheduler0-go-client

echo "Current tags:"
git tag --list | grep v1

echo ""
echo "Creating v1.1.0 tag..."
git tag v1.1.0
git push origin v1.1.0

echo ""
echo "Updated tags:"
git tag --list | grep v1
