FROM golang:1.20 as serverbuild

WORKDIR /app

COPY go.mod go.sum ./

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./totp-service main.go

FROM node:16-alpine AS frontbuild

WORKDIR /app

COPY ./ui/package.json ./

RUN yarn

COPY ./ui .

RUN yarn build

FROM alpine:latest as runner
RUN apk --no-cache add ca-certificates

COPY --from=serverbuild /app/totp-service  .
COPY emqxsl-ca.crt  /emqxsl-ca.crt
COPY data/.gitkeep  /data/
COPY --from=frontbuild /app/dist/  /ui/dist

EXPOSE 8080
# Run
CMD ["/totp-service"] 