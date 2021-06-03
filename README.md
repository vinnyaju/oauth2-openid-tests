# Testes com OAuth2 e OpenID Connect

![https://img.shields.io/github/license/vinnyaju/oauth2-openid-tests?style=plastic](https://img.shields.io/github/license/vinnyaju/oauth2-openid-tests?style=plastic) ![enter image description here](https://img.shields.io/badge/Status-Em%20constru%C3%A7%C3%A3o-orange?style=plastic)

Esse repositório, em construção, contem uma série de scripts que montam um cenário de teste de segurança de APIs, a ideia é utilizar **OAuth2** e **OpenID Connect**, como **IAM (Gerenciamento de Identidades e Acesso)** estou usando o **Keycloack**.

# O que esse repositório faz?

* Cria uma VM Ubuntu no virtualbox com 4 processadores e 4gb de RAM d IP Fixo 192.168.100.101;
* Instala o Guest Additions;
* Adiciona a chave pública configurada no seu github para login na VM criada;
* Atualiza o SO;
* Instala do Docker;
* Instala do docker-compose;
* Faz o deploy de uma imagem do Keycloack com Banco de Dados Mysql.

# Requisitos
* Vagrant;
* Virtualbox;

## O que fazer?

Clone o repositório, coloque o seu usuário Github no valor da variável GithubUser no Vagrantfile e dê um:

    vagrant up

Quando terminar:

    vagrant ssh

Acesse o IP da VM, porta 8080 pelo browser, interface do Keycloack, login / password: admin / Pa55w0rd