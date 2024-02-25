# Use an official Node.js alpine image as a base
FROM node:alpine AS style

# Set the working directory in the container
WORKDIR /app

# Copy the local assets directory to the container
COPY ./assets/ /app/assets
COPY ./tailwind.config.js ./tailwind.config.js

# Install dependencies
RUN npm install -g tailwindcss

# Run the Tailwind CSS build command
RUN npx tailwindcss build ./assets/input.css -o style.css

# Use official golang image as a build base
FROM golang:alpine AS build
WORKDIR /src
COPY . .
COPY --from=style /app/style.css /src/assets/style.css

# Install dependencies
RUN go mod tidy

# Make directories for build dir
RUN mkdir build
RUN mkdir build/files

# Install c dependencies
RUN apk add gcc musl-dev

# Build the application.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage a bind mount to the current directory to avoid having to copy the
# source code into the container.
RUN GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-s -w" -o ./build/server

################################################################################
# Create a new stage for running the application that contains the minimal
# runtime dependencies for the application. This often uses a different base
# image from the build stage where the necessary files are copied from the build
# stage.
#
# The example below uses the alpine image as the foundation for running the app.
# By specifying the "latest" tag, it will also use whatever happens to be the
# most recent version of that image when you build your Dockerfile. If
# reproducability is important, consider using a versioned tag
# (e.g., alpine:3.17.2) or SHA (e.g., alpine@sha256:c41ab5c992deb4fe7e5da09f67a8804a46bd0592bfdf0b1847dde0e0889d2bff).
FROM alpine:latest AS final

# Install any runtime dependencies that are needed to run your application.
# Leverage a cache mount to /var/cache/apk/ to speed up subsequent builds.
RUN apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

# Copy build assets
COPY --from=build /src/build /build
WORKDIR /build

# Create a non-privileged user that the app will run under.
# See https://docs.docker.com/go/dockerfile-user-best-practices/
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

# Set directory owner to unprivalaged user
RUN chown -R appuser:appuser .

# Change user to unprivalaged user
USER appuser

# Expose the port that the application listens on.
EXPOSE 8080

# What the container should run when it is started.
CMD [ "./server" ]
