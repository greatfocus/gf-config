    FROM golang:1.16-alpine3.12 as build
     
    WORKDIR /source
    COPY . .
     
    ARG COMMIT
    RUN CGO_ENABLED=0 go mod tidy
    RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.commit=${COMMIT}" -o bin/gf-config main.go
     
    FROM alpine:3.12
     
    COPY --from=build /source/bin/gf-config /bin/application/gf-config
    COPY --from=build /source/vault /bin/application/vault
     
    EXPOSE 5000
     
    ENTRYPOINT ["./bin/application/gf-config"]