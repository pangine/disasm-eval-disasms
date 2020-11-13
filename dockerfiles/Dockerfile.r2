FROM pangine/disasms-base

# Install disassembler required packages
USER root
WORKDIR /root/
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y \
    curl

USER ${USER}
WORKDIR ${USER_HOME}

# Install r2
RUN wget --progress=bar:force:noscroll https://github.com/radareorg/radare2/archive/4.4.0.tar.gz && \
    tar zxf 4.4.0.tar.gz && \
    rm 4.4.0.tar.gz && \
    cd radare2-4.4.0 && \
    ./configure && \
    make && \
    sys/user.sh

ENV PATH="${USER_HOME}/bin:${PATH}"

# Install this package
RUN go get -u github.com/pangine/disasm-eval-disasms/... && \
    echo "[2020-11-12]"
