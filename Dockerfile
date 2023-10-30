#  Build the Frontend
FROM node:lts-bullseye-slim@sha256:d93fb5c25db163dc795d40eabf66251a2daf6a2c6a2d21cc29930e754aef4c2c  AS frontend
WORKDIR /builder
RUN npm config set unsafe-perm true && npm install -g pnpm
COPY /cmd/frontend/package.json /cmd/frontend/pnpm-lock.yaml ./
RUN pnpm install
COPY /cmd/frontend ./
RUN pnpm build

# Build binary
FROM golang:1.21.1 AS binary
WORKDIR /builder
ENV GOOS=linux GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
COPY --from=frontend /builder/dist ./cmd/frontend/dist
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-w -s -extldflags "-static"' -a -o ./bin/server cmd/api/main.go
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-w -s -extldflags "-static"' -a -o ./bin/migrate cmd/migrate/main.go

# Stage 3: Run the binary
FROM alpine:3.17
RUN apk update && \
  apk add --no-cache ca-certificates openssl dumb-init && \
  rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=binary /builder/bin ./
COPY ./start.sh ./
RUN chmod +x ./start.sh
RUN adduser -D go
RUN chown go:go /app/
EXPOSE ${PORT}
USER go
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/bin/sh", "-c", "./start.sh"]
