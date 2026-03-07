FROM rockylinux/rockylinux:10.1.20251123-minimal

WORKDIR /test

RUN microdnf install -y wget tar
RUN wget https://github.com/astral-sh/python-build-standalone/releases/download/20260303/cpython-3.10.20+20260303-x86_64-unknown-linux-gnu-install_only.tar.gz -O /opt/python-3.10.tar.gz && \
cd /opt/ && \
tar xvf /opt/python-3.10.tar.gz && \
cp -r /opt/python/./ /usr/./ && \
rm -rf /opt/python-3.10.tar.gz /opt/python


ENTRYPOINT ["python3", "/test/e2e/verify.py"]