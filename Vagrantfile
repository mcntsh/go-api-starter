# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "phusion/ubuntu-14.04-amd64"
  config.vm.provision :shell, path: "scripts/provisioning/vagrant-ubuntu.sh"
  config.vm.network :private_network, type: :static, ip: "192.168.50.240"

  config.ssh.insert_key    = true
  config.ssh.forward_agent = true

  config.vm.synced_folder ENV['GOPATH'], "/go"
end
