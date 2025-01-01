FROM alpine:3.21@sha256:21dc6063fd678b478f57c0e13f47560d0ea4eeba26dfc947b2a4f81f686b9f45 AS builder
COPY . .
RUN apk --no-cache add go && CGO_ENABLED=0 go build ./cmd/ccv
FROM alpine:3.21@sha256:21dc6063fd678b478f57c0e13f47560d0ea4eeba26dfc947b2a4f81f686b9f45
RUN apk --no-cache add git
COPY --from=builder ccv /usr/local/bin/
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
