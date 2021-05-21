# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/rurick/coinswallet

# Build the coinswallet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/rurick/coinswallet/cmd/wallet@latest

# Run the coinswallet command by default when the container starts.
ENTRYPOINT /go/bin/wallet

# Document that the service listens on port 8081.
EXPOSE 8081