FROM pangine/disasms-base

# Install disassembler required packages
USER root
WORKDIR /root/
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y \
    curl

USER ${USER}
WORKDIR ${USER_HOME}

# Install r2
RUN curl -Ls https://github.com/radareorg/radare2/releases/download/5.8.8/radare2-5.8.8.tar.xz | tar xJv && \
    radare2-5.8.8/sys/user.sh

ENV PATH="${USER_HOME}/bin:${PATH}"

# Install this package
RUN go get -u github.com/pangine/disasm-eval-disasms/... && \
    go install github.com/pangine/disasm-eval-disasms/... && \
    echo "[2024-01-26]"
