FROM golang:1.22-alpine3.20 AS builder

WORKDIR /usr/src/app

RUN apk add make

# to be removed once there is no more pkg injection.
COPY --from=pkg . /usr/pkg

COPY go.mod go.sum ./

RUN go mod download && go mod verify

CMD ["sh"]

# =======================================================

FROM builder AS app

COPY . .
RUN mkdir -p ./bin
RUN make build

RUN cp ./bin/uda /usr/local/bin/uda

CMD ["uda"]

# =======================================================
FROM builder AS dev

## Copy base makefiles
COPY Makefile ./Makefile

RUN make tools

CMD ["go run main.go"]
