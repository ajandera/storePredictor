FROM golang:1.17-alpine

RUN apk add git
RUN git config --global url."https://azure:glpat-K7tZ-sq9Zxk9Due9d142@gitlab.eaineu.com".insteadOf \
                          "https://gitlab.eaineu.com"

COPY ./entrypoint.sh /
