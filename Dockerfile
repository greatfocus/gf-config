    FROM golang:1.15-alpine3.12 as build
     
    WORKDIR /source
    COPY . .
     
    ARG COMMIT
    RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.commit=${COMMIT}" -o bin/pipeline main.go
     
    FROM alpine:3.12
     
    COPY --from=build /source/bin/pipeline /bin/pipeline
     
    EXPOSE 8080
     
    ENTRYPOINT ["./bin/pipeline"]