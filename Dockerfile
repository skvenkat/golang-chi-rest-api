FROM golang:alpine AS build

RUN apk --no-cache add \
    gcc \
    g++ \
    make \
    git

RUN apk --no-cache add \
    npm \
    yarn
RUN npm install typescript -g

WORKDIR /go/src/app

COPY ./cmd ./cmd
COPY ./configs ./configs
COPY ./internal ./internal

COPY ./web ./web

COPY go.mod .
COPY go.sum .
COPY *.go ./
RUN go mod tidy
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/apiserver ./main.go

RUN cd web && npm install && yarn build && cd ..

FROM alpine:3.17
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /go/src/app/bin /app
COPY --from=build /go/src/app/configs /app/configs

COPY --from=build /go/src/app/web/dist /app/web/dist

EXPOSE 8080

ENTRYPOINT /app/apiserver run --deployment=local
