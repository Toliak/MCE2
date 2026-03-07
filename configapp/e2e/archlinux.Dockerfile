FROM archlinux:base-20260301.0.494762

WORKDIR /test

RUN pacman -Syy && pacman --noconfirm -S wget zstd tar
RUN wget https://github.com/astral-sh/python-build-standalone/releases/download/20260303/cpython-3.10.20+20260303-x86_64-unknown-linux-musl-noopt+static-full.tar.zst -O /opt/python-3.10.tar.zst && \
    cd /opt/ && \
    tar --zstd -xvf /opt/python-3.10.tar.zst && \
    cp -r /opt/python/install/./ /usr/./ && \
    rm -rf /opt/python-3.10.tar.zst /opt/python

ENTRYPOINT ["python3", "/test/e2e/verify.py"]
