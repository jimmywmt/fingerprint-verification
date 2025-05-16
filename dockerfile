FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

# RUN git clone https://github.com/jimmywmt/fingerprint-verification.git .
COPY . .

# 建立輸出目錄
RUN mkdir -p /out

# 安裝必要工具
RUN apk add --no-cache git upx && \
    go install mvdan.cc/garble@latest

ARG XOR_KEY
ARG SHARED_SECRET

# 隨機 XOR 金鑰
RUN set -eux; \
  bytes=$(echo -n "$SHARED_SECRET" | od -An -t u1 | tr -d '\n' | xargs -n1); \
  encoded=""; \
  for b in $bytes; do \
    encoded="${encoded}$((b ^ XOR_KEY)),"; \
  done; \
  encoded="${encoded%,}"; \
  { \
    echo 'package main'; \
    echo 'func decodeXOR(data []byte, key byte) string {'; \
    echo '  dec := make([]byte, len(data))'; \
    echo '  for i := range data { dec[i] = data[i] ^ key }'; \
    echo '  return string(dec)'; \
    echo '}'; \
    echo 'func GetSharedSecret() string {'; \
    echo "  return decodeXOR([]byte{${encoded}}, byte(${XOR_KEY}))"; \
    echo '}'; \
  } > keys.go

RUN garble -literals -tiny -seed=random build -ldflags="-s -w" -o /out/verify .


FROM scratch
COPY --from=builder /out/verify /verify
ENTRYPOINT ["/verify"]
