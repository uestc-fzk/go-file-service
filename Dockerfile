FROM golang:1.16.14
# 此COPY命令将当前目录的所有文件复制到指定目录下
COPY . /opt/GoFileService
WORKDIR /opt/GoFileService
CMD cd /opt/GoFileService \
    && chmod 777 main \
    && ./main
