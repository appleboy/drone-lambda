FROM alpine:3.17

ARG TARGETOS
ARG TARGETARCH

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>" \
  org.label-schema.name="AWS Lambda Plugin" \
  org.label-schema.vendor="Bo-Yi Wu" \
  org.label-schema.schema-version="1.0"

LABEL org.opencontainers.image.source=https://github.com/appleboy/drone-lambda
LABEL org.opencontainers.image.description="Deploying Lambda code with drone CI to an existing function."
LABEL org.opencontainers.image.licenses=MIT

RUN apk add --no-cache ca-certificates=20230506-r0 && \
  rm -rf /var/cache/apk/*

COPY release/${TARGETOS}/${TARGETARCH}/drone-lambda /bin/

ENTRYPOINT ["/bin/drone-lambda"]
