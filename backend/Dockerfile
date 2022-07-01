# Build Stage
# First pull Golang image
FROM golang:1.18.3 as build-env

# Set envirment variable
ENV APP_NAME mz-mock
ENV GO111MODULE=off
# possible values are windows, darwin, js
ENV GOOS linux

# Copy application data into image
COPY . $GOPATH/src/$APP_NAME
WORKDIR $GOPATH/src/$APP_NAME

RUN go get github.com/jackc/pgproto3
# Budild application
RUN CGO_ENABLED=0 go build -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH

# Run Stage
FROM alpine:3.14

# Set envirment variable
ENV APP_NAME mz-mock

# Copy only required data into this image
COPY --from=build-env /$APP_NAME .

# Expose application port
EXPOSE 6875

# Start app 
CMD ./$APP_NAME