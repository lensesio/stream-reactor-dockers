FROM debian

COPY cli-install.sh .

RUN apt-get update && apt-get install -y libssl-dev libffi-dev python-dev build-essential wget curl ssh && \
    wget https://storage.googleapis.com/kubernetes-helm/helm-v2.2.3-linux-amd64.tar.gz && \
    tar zxvf helm-v2.2.3-linux-amd64.tar.gz && \
    cp linux-amd64/helm /usr/local/bin/ && \
    curl -O https://storage.googleapis.com/kubernetes-release/release/v1.6.1/bin/linux/amd64/kubectl && \
    chmod +x ./kubectl && \
    cp kubectl /usr/bin && \
    wget https://github.com/Eneco/landscaper/releases/download/1.0.1/landscaper-v1.0.1-linux-amd64.gz && \
    gunzip landscaper-v1.0.1-linux-amd64.gz && \
    chmod +x landscaper-v1.0.1-linux-amd64 && \
    mv landscaper-v1.0.1-linux-amd64 landscaper && \
    mv landscaper /usr/local/bin/ && \
    wget https://azurecliprod.blob.core.windows.net/install.py && \
    chmod +x install.py && \
    ./cli-install.sh && \
    rm -f install.sh get_helm.sh && \
    echo "source <(kubectl completion bash)" >> ~/.bashrc && \
    wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64 && \
    chmod +x /usr/local/bin/dumb-init && \
    wget https://github.com/Landoop/coyote/releases/download/v1.1/coyote-1.1-linux-amd64 && \
    mv coyote-1.1-linux-amd64 coyote && \
    mv coyote /usr/local/bin
