FROM busybox:1.33.1
WORKDIR /
COPY cmd/web/dist /static
ENTRYPOINT cp -R /static/. /web/
