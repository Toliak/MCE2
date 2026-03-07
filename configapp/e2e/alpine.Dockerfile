FROM alpine:3.23.3

WORKDIR /test

RUN apk add --no-cache wget zstd tar
RUN wget https://github.com/astral-sh/python-build-standalone/releases/download/20260303/cpython-3.10.20+20260303-x86_64-unknown-linux-musl-noopt+static-full.tar.zst -O /opt/python-3.10.tar.zst && \
    cd /opt/ && \
    tar --zstd -xvf /opt/python-3.10.tar.zst

RUN apk add --no-cache sudo shadow
RUN printf "\n%s\n" 'user ALL=(ALL:ALL) NOPASSWD: ALL' >> /etc/sudoers && \
    useradd user -m -u 1000

USER 1000:1000

ENTRYPOINT ["/opt/python/install/bin/python3.10", "/test/e2e/verify.py"]
