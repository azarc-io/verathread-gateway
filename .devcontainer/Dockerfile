# this image is built and updated weekly
# https://github.com/cse-labs/codespaces-images

FROM golang:1.22.4

# some images require specific values
ARG USERNAME=vscode

# log the docker build start event
RUN echo "$(date +'%Y-%m-%d %H:%M:%S')    docker build start" >> /home/$USERNAME/status

RUN chown -R $USERNAME:$USERNAME /home/$USERNAME

# log the docker build complete
RUN echo "$(date +'%Y-%m-%d %H:%M:%S')    docker build complete" >> /home/$USERNAME/status
