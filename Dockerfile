FROM ubuntu:latest
LABEL authors="autowise"

ENTRYPOINT ["top", "-b"]