FROM plugins/base:linux-arm

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>" \
  org.label-schema.name="Drone lambda" \
  org.label-schema.vendor="Bo-Yi Wu" \
  org.label-schema.schema-version="1.0"

COPY release/linux/arm/drone-lambda /bin/

ENTRYPOINT ["/bin/drone-lambda"]
