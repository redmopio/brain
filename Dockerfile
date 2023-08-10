# Specifies a parent image
FROM golang:1.21

# Creates an app directory to hold your app’s source code
WORKDIR /app

# Copies everything from your root directory into /app
COPY . .

# Installs Go dependencies
RUN go mod download

# Builds your app with optional configuration
RUN go build ./cmd/brain

# Create a directory to hold sqlite databases
RUN mkdir -p /app/data

# Tells Docker which network port your container listens on
EXPOSE 8080

# Specifies the executable command that runs when the container starts
CMD [ "/app/brain" ]

