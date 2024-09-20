# Сборка из корневой папки:
# docker build -t todo_list_123 -f build/Dockerfile .

# Запуск:
# docker run -p 7540:7540 todo_list_123

# Запуск с переменными окружения:
# docker run -e TODO_PASSWORD="sample" -p 7540:7540 todo_list_123 

############################

# syntax=docker/dockerfile:1
FROM golang:1.22-alpine3.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /todo_list

EXPOSE 7540

CMD ["/todo_list"]
