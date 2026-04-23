FROM node:22-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
COPY frontend/ .
RUN pnpm build

FROM golang:1.25-alpine AS backend-builder
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/build ./frontend/build
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o litepay .

FROM golang:1.25-alpine AS atlas-builder
RUN go install ariga.io/atlas/cmd/atlas@latest

FROM alpine:latest
RUN apk add --no-cache ca-certificates netcat-openbsd
COPY --from=backend-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend-builder /app/litepay /litepay
COPY --from=atlas-builder /go/bin/atlas /atlas
COPY --from=backend-builder /app/atlas.hcl /atlas.hcl
COPY --from=backend-builder /app/migrations /migrations
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
USER 65534:65534
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
