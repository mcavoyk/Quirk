FROM scratch
ADD bin/quirk ./
ADD api/config.toml ./
EXPOSE 5005
ENTRYPOINT ["/quirk"]