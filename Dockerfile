FROM ubuntu:latest
LABEL authors="christianrodriguez"

ENTRYPOINT ["top", "-b"]