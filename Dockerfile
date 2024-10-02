# Etapa de construcción
FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Construir el binario, deshabilitando CGO y optimizando para tamaño
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sirca-backend ./cmd/server

# Etapa final
FROM alpine:3.18

# Instalar certificados y limpiar caché en la misma capa para reducir tamaño
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copia el binario desde la etapa de construcción
COPY --from=builder /app/sirca-backend .

EXPOSE 8080

CMD ["./sirca-backend"]
