# # # # # # # #
# Build Stage
FROM golang:1.16 AS builder

ENV GO111MODULE=on
ARG SERVICE_NAME

WORKDIR $GOPATH/src/bitbucket.org/bridce/ms-pari-web

ADD . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o app cmd/api/main.go

# Final Stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta

WORKDIR /root/

COPY --from=builder /go/src/bitbucket.org/bridce/ms-pari-web/app .
COPY --from=builder /go/src/bitbucket.org/bridce/ms-pari-web/.env.example ./.env
COPY --from=builder /go/src/bitbucket.org/bridce/ms-pari-web/internal/pkg/config/rbac_model.conf ./internal/pkg/config/

CMD ["./app"]
EXPOSE 4000 8080