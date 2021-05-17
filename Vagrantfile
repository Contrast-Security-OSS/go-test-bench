# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  config.vm.define "centos" do |centos|
    centos.vm.box = "centos/8"
    centos.vm.box_version = "1905.1"
    centos.vm.network "forwarded_port", guest: 8080, host: 18080
    centos.vm.provision :shell, path: "provision/centos.sh"
  end

  config.vm.define "ubuntu" do |ubuntu|
    ubuntu.vm.box = "bento/ubuntu-20.04"
    ubuntu.vm.network "forwarded_port", guest: 8080, host: 28080
    ubuntu.vm.provision :shell, path: "provision/ubuntu.sh"
  end
end
