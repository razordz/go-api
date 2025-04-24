FROM golang:1.23-bookworm

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates git && update-ca-certificates

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# ✅ Agora o docs já existe localmente, só copiamos
RUN go build -o main .

EXPOSE 8080
CMD ["./main"]
