FROM centos

COPY entrypoint.sh entrypoint.sh

RUN yum install -y wget && \
    wget https://archive.cloudera.com/cdh5/redhat/7/x86_64/cdh/cloudera-cdh5.repo?_ga=1.26965246.2135856916.1478605511 \
    -O /etc/yum.repos.d/cloudera.repo && \
    wget http://archive.cloudera.com/beta/kudu/redhat/7/x86_64/kudu/cloudera-kudu.repo \
    -O /etc/yum.repos.d/kudu.repo && \
    yum install -y gcc gcc-c++ autoconf automake \
    epel-release python-pip impala-shell hive hbase spark-core spark-python \
    postgresql-devel python-devel java-1.8.0-openjdk.x86_64 unzip kudu-client-devel && \
    wget https://github.com/edenhill/librdkafka/archive/v0.9.2.zip -O v0.9.2.zip && \
    unzip v0.9.2.zip && \
    cd librdkafka-0.9.2/ && \
    ./configure --prefix=/usr && \
    make -j && \
    make install && \
    easy_install -U setuptools && \
    easy_install pip && \
    pip install ibis-framework && \
    pip install "confluent-kafka[avro]" && \
    wget http://packages.confluent.io/archive/3.0.1/confluent-3.0.1-2.11.tar.gz -O /etc/confluent/ && \
    cd .. && \
    chmod +x entrypoint.sh 

ENV HADOOP_CONF_DIR=/etc/yarn-conf
ENV HIVE_CONF_DIR=/etc/hive-conf

CMD ["./entrypoint.sh"]
