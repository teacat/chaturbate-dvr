FROM golang:latest

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build

CMD [ "sh", "-c", "./chaturbate-dvr -u $USERNAME" ]