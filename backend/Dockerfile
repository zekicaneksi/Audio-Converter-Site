FROM golang:alpine

# Set destination for COPY
WORKDIR /app

# Install ffprobe
RUN apk add ffmpeg

# Download Go modules
COPY go.mod ./
RUN go mod download

COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /built

EXPOSE 8080

# Run
CMD ["/built"]