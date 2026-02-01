FROM alpine:3.23@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659 AS builder
COPY . .
RUN apk --no-cache add go && CGO_ENABLED=0 go build ./cmd/ccv
FROM alpine:3.23@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659
RUN apk --no-cache add git
COPY --from=builder ccv /usr/local/bin/
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
