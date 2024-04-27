FROM golang:1.21 AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /main .

FROM golang:1.21 AS final

WORKDIR /usr/src/app

COPY --from=build /usr/src/app/ .

EXPOSE 8080

CMD ["go", "run", "main.go"]
