FROM ubuntu:latest
LABEL authors="jinagyu"

ENTRYPOINT ["top", "-b"]