FROM golang:1.17.4-alpine3.14 as builder

ENV BASE_APP_DIR /go/src/github.com/kyma-incubator/compass/components/system-broker
WORKDIR ${BASE_APP_DIR}

COPY go.mod go.sum ${BASE_APP_DIR}/
RUN go mod download -x

COPY . ${BASE_APP_DIR}

RUN go build -v -o /app/system-broker ./cmd/main.go

FROM alpine:3.14.2
LABEL source = git@github.com:kyma-incubator/compass.git
WORKDIR /app

RUN apk --no-cache add curl ca-certificates

COPY --from=builder /app /app

CMD ["/app/system-broker"]
