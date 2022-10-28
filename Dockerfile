from golang:latest

WORKDIR /welsh-academy
COPY ./ ./

ENV PORT=9000
ENV DB_HOST="localhost"
ENV DB_PORT=5432
ENV DB_USER="admin"
ENV DB_PASS="admin"
ENV DB_NAME="welsh"

RUN go mod download

RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN task gen-swagger
RUN task build

ENTRYPOINT ["./bin/welsh-academy"]
