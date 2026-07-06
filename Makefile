# Scheduler0 Go client — release automation
#
# Usage:
#   make test                  Run tests
#   make release VERSION=1.1.3 Tag vX.Y.Z, push, and warm the module proxy
#
# Go modules are not "published" to a registry: a release is simply a semver
# git tag. The first `go get ...@vX.Y.Z` makes the public proxy fetch it.

MODULE      := github.com/scheduler0/scheduler0-go-client
MAIN_BRANCH := main

.PHONY: help test vet tidy release \
        guard-VERSION check-clean check-branch check-tag

help:
	@echo "make test                  - run tests"
	@echo "make release VERSION=1.1.3 - tag vX.Y.Z, push, warm the module proxy"

test:
	go test ./...

vet:
	go vet ./...

tidy:
	go mod tidy

release: guard-VERSION check-branch check-clean check-tag vet test
	@echo ">> Releasing $(MODULE) v$(VERSION)"
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)
	@echo ">> Warming module proxy (first fetch may take a moment)..."
	GOPROXY=proxy.golang.org GOFLAGS= go list -m $(MODULE)@v$(VERSION) || true
	@echo ">> Released $(MODULE) v$(VERSION)"

# --- guards -----------------------------------------------------------------

guard-VERSION:
	@if [ -z "$(VERSION)" ]; then echo "ERROR: VERSION is required, e.g. make release VERSION=1.1.3"; exit 1; fi
	@echo "$(VERSION)" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+(-.*)?$$' || { echo "ERROR: VERSION '$(VERSION)' is not a valid semver (expected x.y.z)"; exit 1; }

check-branch:
	@if [ "$$(git branch --show-current)" != "$(MAIN_BRANCH)" ]; then echo "ERROR: not on '$(MAIN_BRANCH)' branch"; exit 1; fi

check-clean:
	@if [ -n "$$(git status --porcelain)" ]; then echo "ERROR: working tree is dirty; commit or stash first"; exit 1; fi

check-tag:
	@if git rev-parse -q --verify "refs/tags/v$(VERSION)" >/dev/null; then echo "ERROR: tag v$(VERSION) already exists"; exit 1; fi
