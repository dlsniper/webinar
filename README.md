## Webinar

This project is just a tech demo for debugging Go apps using Gogland and containers.

### Requirements

This needs Docker 17.05 or newer in order to support multi-stage building. For more information
please see the manual: https://docs.docker.com/engine/userguide/eng-image/multistage-build/

### Usage

#### Building

To build the container use:
```bash
docker build -t webinar:debug 
```

#### Docker < 17.05

If a newer version of Docker is not available, you can still use ` container-debug-old.sh `.


#### Running

To run the container use:

```bash
docker run --rm \
    --name=webinar-debug \
    -p 8000:8000 \
    -p 2345:40000 \
    --security-opt=seccomp:unconfined \
    webinar:debug
```

### What does it do?

This builds a Go binary then it adds Delve to a Docker container and runs them both using:
 
 ```bash
/dlv --listen=:40000 --headless=true --api-version=2 exec /webinar
```

### License

This project is under Apache 2.0 license, please see the LICENSE file.
