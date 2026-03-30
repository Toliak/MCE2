FROM debian:trixie-20260223

WORKDIR /test

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update -y --allow-releaseinfo-change && \
    apt-get install -y wget tar
RUN wget https://github.com/astral-sh/python-build-standalone/releases/download/20260303/cpython-3.10.20+20260303-x86_64-unknown-linux-gnu-install_only.tar.gz -O /opt/python-3.10.tar.gz && \
    cd /opt/ && \
    tar xvf /opt/python-3.10.tar.gz

RUN apt-get install -y sudo
RUN printf "\n%s\n" 'user ALL=(ALL:ALL) NOPASSWD: ALL' >> /etc/sudoers && \
    useradd user -m -u 1000

USER 1000:1000

ENTRYPOINT ["/opt/python/bin/python3.10", "/test/e2e/overlay/__main__.py"]
