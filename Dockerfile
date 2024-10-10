FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o /go_final_project

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/web /app/web
COPY --from=build /go_final_project /app/go_final_project
ENV TODO_PORT=8080
ENV TODO_DBFILE=scheduler.db
EXPOSE ${TODO_PORT}
CMD ["/app/go_final_project"]