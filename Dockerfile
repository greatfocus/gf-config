    FROM golang:1.15-alpine3.12 as build
     
    WORKDIR /source
    COPY . .
     
    ARG COMMIT
    RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.commit=${COMMIT}" -o bin/gf-config main.go
     
    FROM alpine:3.12
     
    COPY --from=build /source/bin/gf-config /bin/gf-config
     
    EXPOSE 8080
     
    ENTRYPOINT ["./bin/gf-config"]