FROM golang:1.20 as builder
WORKDIR /airgap-webhook
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags='-w -s' -o /go/bin/airgap-webhook .

# Runtime Image
FROM scratch
COPY --from=builder /go/bin/airgap-webhook /bin/airgap-webhook
ENTRYPOINT [ "/bin/airgap-webhook" ]
