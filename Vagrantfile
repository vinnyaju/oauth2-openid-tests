# -*- mode: ruby -*-
# vi: set ft=ruby :
required_plugins = %w(vagrant-vbguest)
required_plugins.each do |plugin|
  system "vagrant plugin install #{plugin}" unless Vagrant.has_plugin? plugin
end

BoxUser = "vagrant"
GithubUser = "vinnyaju"

Vagrant.configure("2") do |config|

	config.vm.define "keycloaksandbox", primary: true do |keycloaksandbox|
		keycloaksandbox.vm.box = "ubuntu/hirsute64"
		keycloaksandbox.vm.hostname = "keycloak"
		keycloaksandbox.vm.network "private_network", ip: "192.168.100.101"

    ##Inserindo minha chave pública do github no authorized_keys da VM (Troque o nome do usuário para o seu para usar as suas chaves públicas)
    keycloaksandbox.vm.provision "shell", path: "shell-scripts/set-key-login-putty.sh", args: [GithubUser, BoxUser]

    ## Atualizando a máquina
    keycloaksandbox.vm.provision "shell", path: "shell-scripts/update-box.sh"

    ## Instalando pacotes utilitários
    keycloaksandbox.vm.provision "shell", path: "shell-scripts/install-aditional-packages.sh"

    ## Instalando o Docker e docker-compose
    keycloaksandbox.vm.provision "shell", path: "shell-scripts/install-docker-engine.sh", args: BoxUser

		## Subindo o docker-compose
    keycloaksandbox.vm.provision "shell", path: "shell-scripts/start-app.sh"

  	## Printando o IP host-only
		keycloaksandbox.vm.provision "shell", inline: "ifconfig enp0s8 | grep inet"

		keycloaksandbox.vm.provider "virtualbox" do |vb|
			vb.gui = false
			vb.name = "keycloak"
			vb.memory = "4096"
			vb.cpus = "4"
		end

	end

  if Vagrant.has_plugin?("vagrant-vbguest")
    config.vbguest.auto_update = true
    config.vbguest.no_remote = true
  end

end