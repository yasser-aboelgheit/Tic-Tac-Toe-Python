FROM golang:1.22-alpine3.20 AS builder

WORKDIR /usr/src/app

RUN apk add make

COPY go.mod go.sum ./

RUN go mod download && go mod verify

CMD ["sh"]

# =======================================================

FROM builder AS app

COPY . .
RUN mkdir -p ./bin
RUN make build

RUN cp ./bin/startupbuilder /usr/local/bin/startupbuilder

CMD ["startupbuilder"]

# =======================================================
FROM builder AS dev

## Copy base makefiles
COPY ./tools/makefiles ./tools/makefiles
COPY Makefile ./Makefile

RUN make install-linter

CMD ["go run main.go"]



