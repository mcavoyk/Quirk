FROM scratch
ADD bin/main ./
ADD api/config.toml ./
EXPOSE 5005
ENTRYPOINT ["/main"]