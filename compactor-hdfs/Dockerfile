FROM centos

COPY entrypoint.sh entrypoint.sh

RUN yum install -y wget && \
    wget https://archive.cloudera.com/cdh5/redhat/7/x86_64/cdh/cloudera-cdh5.repo?_ga=1.26965246.2135856916.1478605511 \
    -O /etc/yum.repos.d/cloudera.repo && \
    yum clean all && \
    yum install -y impala-shell hive java-1.8.0-openjdk.x86_64 unzip && \
    chmod +x entrypoint.sh

ENV HADOOP_CONF_DIR=/etc/yarn-conf
ENV HIVE_CONF_DIR=/etc/hive-conf

CMD ["./entrypoint.sh"]
