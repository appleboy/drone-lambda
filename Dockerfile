FROM plugins/base:linux-amd64

COPY release/linux/amd64/drone-lambda /bin/

ENTRYPOINT ["/bin/drone-lambda"]
