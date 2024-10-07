FROM golang:1.22 AS build

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go_final_project

FROM ubuntu:latest
WORKDIR /app
COPY --from=build /app/web /app/web
COPY --from=build /go_final_project /app/go_final_project
ENV TODO_PORT=8080
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=qwerty
EXPOSE 8080
CMD ["/app/go_final_project"]