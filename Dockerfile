FROM scratch
COPY bin/quirk_api api/config.toml ./
EXPOSE 5005
ENTRYPOINT ["/quirk_api"]