FROM node:22-bookworm AS frontend-builder
WORKDIR /app
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
COPY frontend/ .
RUN pnpm build

FROM golang:1.25-bookworm AS backend-builder
RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Download Atlas CLI binary directly
RUN curl -sSfL https://atlasbinaries.com/atlas/atlas-linux-amd64-latest -o /usr/local/bin/atlas && chmod +x /usr/local/bin/atlas
COPY . .
COPY --from=frontend-builder /app/build ./frontend/build
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o litepay .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates netcat-openbsd && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /go /home/nonroot && chown -R 65534:65534 /go /home/nonroot
COPY --from=backend-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend-builder /usr/local/bin/atlas /usr/local/bin/atlas
COPY --from=backend-builder /app/atlas.hcl /atlas.hcl
COPY --from=backend-builder /app/migrations /migrations
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
USER 65534:65534
ENV HOME=/home/nonroot
ENV GOCACHE=/go
ENV ATLAS_CACHE=/home/nonroot/.atlas/cache
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
