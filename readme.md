# navget-cli

This is a command-line tool for creating Navget packages and uploading them.

It also provides a Dockerfile for a Drone plugin.

### How to update the Drone plugin

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o navget-cli
docker build --tag docker.naveego.com:4333/navget-drone:latest .
docker push docker.naveego.com:4333/navget-drone:latest
```

Configuration of the Drone plugin in a .drone.yml file looks like this:

```
  publish-test:
    image: docker.naveego.com:4333/navget-drone:latest
    pull: true
    # Navget endpoint, with no trailing /
    endpoint: http://navget.naveego.test
    # Space-delimited list of files to include in the package
    files: plugin-pub-test     
    # JWT token for Navget (or provide a secret named NAVGET_TOKEN)
    token: {... token ... }
    # The OS to specify in the Navget upload
    os: ${GOOS}

```

### How to run the tool

The tool is self-documenting (run `navget-cli --help` to see help).

All tool parameters can set through environment variables by uppercasing and prefixing them with NAVGET_.