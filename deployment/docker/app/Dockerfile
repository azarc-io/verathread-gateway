FROM node:20-alpine
LABEL authors="waeljammal"

ARG BINARY

ENV BINARY=${BINARY}

COPY $BINARY /

ENTRYPOINT ["/vth-gateway"]
