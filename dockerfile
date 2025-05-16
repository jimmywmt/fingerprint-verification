FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

RUN git clone https://github.com/jimmywmt/fingerprint-verification.git .

# 建立輸出目錄
RUN mkdir -p /out

# 安裝必要工具
RUN apk add --no-cache git upx && \
    go install mvdan.cc/garble@latest

ARG XOR_KEY
ARG SHARED_SECRET

# 根據實際情況調整這一行路徑
# RUN sed -i "s|{{SHARED_SECRET}}|${SHARED_SECRET}|g" main.go && \
#     go build -ldflags="-s -w" -o /out/verify .

# 隨機 XOR 金鑰
RUN set -eux; \
  enc() { \
    key=$1; str=$2; \
    bytes=$(echo -n "$str" | od -An -t u1 | tr -d '\n' | xargs -n1); \
    encoded=""; \
    for b in $bytes; do \
      encoded="${encoded}$(($b ^ $XOR_KEY)),"; \
    done; \
    encoded="${encoded%,}"; \
    echo "var $key = decodeXOR([]byte{${encoded}}, $XOR_KEY)" >> keys.go; \
  }; \
  echo 'package main' > keys.go; \
  echo 'func decodeXOR(data []byte, key byte) string { dec := make([]byte, len(data)); for i := range data { dec[i] = data[i] ^ key }; return string(dec) }' >> keys.go; \
  enc SharedSecret "${SHARED_SECRET}" \
  garble -literals -tiny -seed=random build -ldflags="-s -w" -o /out/verify .


FROM scratch
COPY --from=builder /out/verify /verify
ENTRYPOINT ["/verify"]
