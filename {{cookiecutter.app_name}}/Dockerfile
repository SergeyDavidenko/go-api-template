FROM golang:{{cookiecutter.docker_build_image_version}} as builder

# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Run test
RUN go test ./...
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go 
#------------------------------------------------------
######## Start a new stage from scratch #######
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/main .
COPY --from=builder /go/src/github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/configs/{{cookiecutter.app_name}}.yml ./configs/{{cookiecutter.app_name}}.yml

# Expose port 8080 to the outside world
EXPOSE 8080
# Command to run the executable
CMD ["./main"] 