#  Build the Frontend
FROM node:lts-bullseye-slim@sha256:d93fb5c25db163dc795d40eabf66251a2daf6a2c6a2d21cc29930e754aef4c2c  AS frontend
WORKDIR /builder
RUN npm config set unsafe-perm true && npm install -g pnpm
COPY /www/fiber-bot/package.json /www/fiber-bot/pnpm-lock.yaml ./
RUN pnpm install
COPY /www/fiber-bot ./
RUN pnpm build

# Build go binary
FROM golang:1.21.1 AS api-build
WORKDIR /builder
ENV GOOS=linux GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
COPY --from=frontend /builder/dist ./www/fiber-bot/dist
RUN go build -o ./bin/server .

FROM alpine:3.17

RUN apk update && \
  apk add --no-cache ca-certificates openssl dumb-init && \
  rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=api-build /builder/bin ./
COPY ./start.sh ./

RUN chmod +x ./server && chmod +x ./start.sh

RUN adduser -D go

RUN chown go:go /app/

EXPOSE ${PORT}

USER go

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/bin/sh", "-c", "./start.sh"]
