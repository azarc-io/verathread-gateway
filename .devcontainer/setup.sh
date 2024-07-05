# install curl, git, ...
apt-get update
apt-get install -y \
  curl \
  git \
  gnupg2 \
  jq \
  sudo \
  zsh

useradd -m user
su user

# install go
VERSION='1.22.2'
OS='linux'
ARCH='amd64'

curl -OL https://dl.google.com/go/go${VERSION}.${OS}-${ARCH}.tar.gz
tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
rm go$VERSION.$OS-$ARCH.tar.gz

INSTALLED_GO_VERSION=$(go version)
echo "Go version ${INSTALLED_GO_VERSION} is installed"

# install gopls, dlv, hey
echo "Getting development tools"
go install golang.org/x/tools/gopls
go install github.com/go-delve/delve/cmd/dlv
go install github.com/rakyll/hey

# vscode-go dependencies
echo "Getting dependencies for the vscode-go plugin "
# via: https://github.com/microsoft/vscode-go/blob/master/.travis.yml
go install github.com/acroca/go-symbols
go install github.com/cweill/gotests/...
go install github.com/davidrjenni/reftools/cmd/fillstruct
go install github.com/haya14busa/goplay/cmd/goplay
go install github.com/mdempsky/gocode
go install github.com/ramya-rao-a/go-outline
go install github.com/rogpeppe/godef
go install github.com/sqs/goreturns
go install github.com/uudashr/gopkgs/cmd/gopkgs
go install github.com/zmb3/gogetdoc
go install golang.org/x/lint/golint
go install golang.org/x/tools/cmd/gorename

# set-up and install yarn
curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
apt-get update && apt-get install yarn -y

# instsal nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

# setup and install oh-my-zsh
sh -c "$(curl -fsSL https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
cp -R /root/.oh-my-zsh /home/$USERNAME
cp /root/.zshrc /home/$USERNAME
sed -i -e "s/\/root\/.oh-my-zsh/\/home\/$USERNAME\/.oh-my-zsh/g" /home/$USERNAME/.zshrc
chown -R $USER_UID:$USER_GID /home/$USERNAME/.oh-my-zsh /home/$USERNAME/.zshrc
