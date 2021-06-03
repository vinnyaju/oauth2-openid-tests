# Testes com OAuth2 e OpenID Connect

![https://img.shields.io/github/license/vinnyaju/oauth2-openid-tests?style=plastic](https://img.shields.io/github/license/vinnyaju/oauth2-openid-tests?style=plastic) ![enter image description here](https://img.shields.io/badge/Status-Em%20constru%C3%A7%C3%A3o-orange?style=plastic)

Esse repositório, em construção, contem uma série de scripts que montam um cenário de teste de segurança de APIs, a ideia é utilizar **OAuth2** e **OpenID Connect**, como **IAM (Gerenciamento de Identidades e Acesso)** estou usando o **Keycloack**.

# Requisitos
* Vagrant;
* Virtualbox;

## O que fazer?

Clone o repositório, coloque o seu usuário Github no valor da variável GithubUser no Vagrantfile e dê um:

    vagrant up

Quando terminar:

    vagrant ssh

Acesse o IP da VM, porta 8080 pelo browser, interface do Keycloack, login / password: admin / Pa55w0rd