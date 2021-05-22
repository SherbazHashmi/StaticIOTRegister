## Running To Stages: Buuild and Run

# golang base image
FROM golang:alpine as builder
## Build

# Set maintainer
LABEL maintainer="Sherbaz Hashmi <sherbaz.hashmi@gmail.com"

# Add git
RUN apk update && apk add --no-cache git


WORKDIR /app

# copy module files

COPY go.mod go.sum ./

# download dependencies

RUN go mod download

# copy source from current dir to working dir inside the container

COPY . .

# build the app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


FROM alpine:latest
## COPY, CONFIGURE & RUN
RUN apk --no-cache add ca-certificates

WORKDIR /root/
# Copying prebuilt binary file from the previous stage.
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Setup Network
EXPOSE 8080

# Running it

CMD ["./main"]