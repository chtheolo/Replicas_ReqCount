FROM golang:latest

# Create the directory for our service
WORKDIR /usr/src/Replicas_ReqCount

# Copy the module files
COPY go.mod .
COPY go.sum .

# Install modules
RUN go mod download

# Copy the source code
COPY . .

# Compile our application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "req_counter_service" main.go
# RUN go build -o /docker-gs-ping

EXPOSE 8083

CMD [ "./req_counter_service" ]