FROM ubuntu:latest

# update apt
RUN apt-get update -y && \
    apt-get -y install --no-install-recommends \
    gnupg \
    bash-completion \
    openssl \
    libssl-dev \
    ca-certificates \
    vim \
    sudo \
    iptables \
    wget \
    curl && \
    rm -rf /var/lib/apt/lists/*
# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/ShangHai /etc/localtime #经测试，不加这一行有时会不生效。或系统重启后也会恢复成UTC时间
RUN echo "Asia/Shanghai" > /etc/timezone
#golang
RUN wget https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.13.4.linux-amd64.tar.gz && \
    rm go1.13.4.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV e=.env.example
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" "$GOPATH/src/code" && chmod -R 777 "$GOPATH"

#kubectl
# RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl && \
#     chmod +x ./kubectl && \
#     mv ./kubectl /usr/local/bin/kubectl && \
#     mkdir /root/.kube && \
#     touch /root/.kube/config

#telepresence
# RUN curl -s https://packagecloud.io/install/repositories/datawireio/telepresence/script.deb.sh | bash && \
#     apt-get install --no-install-recommends telepresence -y

WORKDIR "$GOPATH/src/code"

# ENV SCOUT_DISABLE=1
