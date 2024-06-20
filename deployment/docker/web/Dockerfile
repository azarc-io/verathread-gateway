FROM busybox
WORKDIR /
COPY cmd/web/dist /static
ENTRYPOINT cp -R /static/. /web/
