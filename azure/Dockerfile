FROM debian

COPY cli-install.sh .
COPY configure.sh .

RUN apt-get update && apt-get install -y libssl-dev libffi-dev python-dev build-essential wget curl ssh && \
    curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh && \
    chmod +x get_helm.sh && \
    sed -i "s/sudo//g" get_helm.sh && \
    ./get_helm.sh && \
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
    echo "source <(kubectl completion bash)" >> ~/.bashrc