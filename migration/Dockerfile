# syntax=docker/dockerfile:1

FROM alpine:3.21.3 

WORKDIR /app

RUN apk add --no-cache bash curl

RUN curl -sSf https://atlasgo.sh | sh

COPY scripts/wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh
