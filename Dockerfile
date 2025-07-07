FROM golang:1.24.4-alpine3.21

RUN apk add --no-cache tzdata curl

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies
# only redownloading them in subsequent builds if they change
COPY go.* ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/ ./...

EXPOSE 8080

COPY healthcheck.sh /usr/local/bin
RUN chmod +x /usr/local/bin/healthcheck.sh

HEALTHCHECK CMD healthcheck.sh

CMD ["segment"]
