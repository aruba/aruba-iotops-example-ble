FROM golang:1.22-alpine AS build

WORKDIR /go/src/app

COPY . .
RUN go build -ldflags '-w -s' -o app ./application

FROM alpine

RUN mkdir -p /home/app

COPY --from=build /go/src/app/app /home/app/app

# Only /home/app/data directory is writable. Read developer documentation for details
CMD /home/app/app
