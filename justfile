APP      := "clippy"
VERSION  := `perl -nE'm{Version\s*=\s*"(\d+\.\d+.\d+)"} && print $1' ./cmd/root.go`
REGISTRY := "gcr.io/docker-registry-mg/mg"
BINARY   := "clippy-linux-amd64"

build:
  echo "Building version {{VERSION}} of {{APP}}"
  go build -o clippy main.go

lint:
  go vet ./... || true
  golangci-lint run ./... || true
  govulncheck ./...

test:
  go test ./...

cross-build:
  echo "Cross-compiling {{APP}} {{VERSION}} for linux/amd64"
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o {{BINARY}} main.go

image-build: cross-build
  podman build -f ci/Dockerfile --platform=linux/amd64 -t mg/{{APP}} .
  rm {{BINARY}}

image-tag:
  podman tag mg/{{APP}}:latest mg/{{APP}}:{{VERSION}}
  podman tag mg/{{APP}}:latest {{REGISTRY}}/{{APP}}
  podman tag mg/{{APP}}:latest {{REGISTRY}}/{{APP}}:{{VERSION}}

image-push:
  podman push {{REGISTRY}}/{{APP}}:{{VERSION}}
  podman push {{REGISTRY}}/{{APP}}:latest

deploy:
  git diff --exit-code
  go vet ./...
  go test ./...
  just image-build
  just image-tag
  git tag "{{VERSION}}"
  git push
  git release
  git push --tags
  just image-push

release:
  git diff --exit-code
  git tag "{{VERSION}}"
  git push
  git release
  git push --tags
  goreleaser release --clean
