APP      := "clip"
# VERSION  := `perl -nE'm{version\s*=\s*"(\d+\.\d+.\d+)"} && print $1' ./cmd/root.go`
VERSION  := "0.0.1"

build:
  echo "Building verions {{VERSION}} of {{APP}}"
  go build -o clip main.go

lint:
  go vet ./... || true
  golangci-lint run ./... || true
  govulncheck ./...

release:
  git diff --exit-code
  git tag "{{VERSION}}"
  git push
  git release
  git push --tags
  goreleaser release --clean
