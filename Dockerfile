FROM alpine:3.20
RUN apk --no-cache add curl git jq \
    && (if echo "$TARGETPLATFORM" | grep -q arm; then \
    curl -sSL $(curl -s https://api.github.com/repos/smlx/ccv/releases/latest | jq -r '.assets[].browser_download_url | select(test("linux_arm64"))'); \
    else \
    curl -sSL $(curl -s https://api.github.com/repos/smlx/ccv/releases/latest | jq -r '.assets[].browser_download_url | select(test("linux_amd64"))'); \
    fi) \
    | tar -xz -C /usr/local/bin ccv
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
