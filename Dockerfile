# You must have built navget-cli using `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o navget-cli` before building this image.
# This image is intended for use a drone plugin named docker.naveego.com:4333/navget-drone:latest
FROM alpine
ADD navget-cli /bin/
RUN apk -Uuv add ca-certificates

# ENV NAVGET_ENDPOINT=http://navget.naveego.test
# ENV NAVGET_TOKEN=
# ENV NAVGET_FILES=
# ENV NAVGET_OS=
# ENV NAVGET_ARCH=

ENTRYPOINT /bin/navget-cli publish