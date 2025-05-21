FROM golang:1.24 AS be-builder

WORKDIR /build

COPY backend .
RUN go mod download

RUN CGO_ENABLED=0 \
    go build \
    -ldflags "-s -w" \
    -trimpath \
    -o server \
    .

FROM node:23-slim AS fe-builder

WORKDIR /build

COPY ./frontend .

RUN npm install
RUN npm run build

# Resulting container image
FROM node:23-slim

WORKDIR /app
COPY --from=fe-builder /build .
COPY --from=be-builder /build/server /app/server
COPY entrypoint.sh /app/entrypoint.sh

RUN mkdir -p /common
RUN mkdir -p /offline
COPY common/ietf-inet-types.yang common/ietf-yang-types.yang /common/
COPY common/nsp-model-extensions.yang common/webfwk-ui-metadata.yang /common/
COPY common/nsp-lso-manager.yang common/nsp-lso-operation.yang /common/

ENV HOST=0.0.0.0
EXPOSE 4173

CMD [ "/app/entrypoint.sh" ]