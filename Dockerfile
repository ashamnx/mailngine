# Stage 1: Build Go binaries
FROM golang:1.26-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/worker ./cmd/worker

# Stage 2: Build Angular frontend
FROM node:20-alpine AS ng-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ .
RUN npx ng build --configuration=production

# Stage 3: Production runtime
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=go-builder /app/bin/server /app/bin/worker ./
COPY --from=ng-builder /app/web/dist/web/browser ./static/
EXPOSE 8080
CMD ["./server"]
