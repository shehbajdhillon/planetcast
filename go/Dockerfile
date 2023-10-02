FROM golang:1.21.1-bullseye

RUN apt-get update -qq && apt-get install ffmpeg -y

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /planetcast-gql-server server.go

EXPOSE 8080
EXPOSE 5432

ENTRYPOINT [ "/planetcast-gql-server" ]

