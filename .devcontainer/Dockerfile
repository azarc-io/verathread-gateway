# this image is built and updated weekly
# https://github.com/cse-labs/codespaces-images

FROM golang:1.22.2

COPY . .

ENV PATH "$PATH:/usr/local/go/bin:/home/user/go/bin:/home/user/.cargo/bin"
RUN bash setup.sh

RUN echo 'export NVM_DIR="$HOME/.nvm"' >> "$HOME/.zshrc"
RUN echo '\n' >> "$HOME/.zshrc"
RUN echo '[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"  # This loads nvm' >> "$HOME/.zshrc"

USER user
