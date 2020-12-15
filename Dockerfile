# Build the manager binary
FROM golang:1.13 as builder

WORKDIR /workspace

# Copy the go source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o arm -mod vendor main.go 

# Using alpine
FROM alpine
WORKDIR /home/spinnaker
COPY --from=builder /workspace/arm /bin/arm
COPY --from=builder /workspace/examples ./examples
COPY --from=builder /workspace/README.md /workspace/arm ./