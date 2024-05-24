FROM golang:1.19-alpine

RUN apk add git
RUN apk add openssh
RUN git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"

COPY ./entrypoint.sh /
