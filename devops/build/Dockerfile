ARG GO_VERSION=1.23.1

FROM docker-x5m-prod.x5.ru/base_image/go:$GO_VERSION-bookworm AS builder

WORKDIR /app

ARG CI_JOB_LOGIN=gitlab-ci-token
ARG CI_JOB_TOKEN

RUN echo "machine scm.x5.ru" > ~/.netrc && \
    echo "  login $CI_JOB_LOGIN" >> ~/.netrc && \
    echo "  password $CI_JOB_TOKEN" >> ~/.netrc

ARG GOPRIVATE=scm.x5.ru/x5m/go-backend/packages/*
ARG GOPROXY
ARG GOSUMDB

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o template -ldflags "-s -w" cmd/server/main.go


FROM docker-x5m-prod.x5.ru/base_image/debian:12-slim AS runtime

WORKDIR /app

COPY --from=builder /app/template .

USER 65535

EXPOSE 8080

ENTRYPOINT [ "/app/template", "server" ]
