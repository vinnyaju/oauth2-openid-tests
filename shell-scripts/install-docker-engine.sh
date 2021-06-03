#!/bin/bash

curl -fsSL https://get.docker.com -o docker-install.sh
sh docker-install.sh
usermod -aG docker $1

sudo curl -L https://github.com/docker/compose/releases/download/1.29.1/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose