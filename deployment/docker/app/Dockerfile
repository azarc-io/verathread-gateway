FROM scratch
LABEL authors="waeljammal"

ARG BINARY

ENV BINARY=${BINARY}

COPY $BINARY /

ENTRYPOINT ["/vth-gateway"]
