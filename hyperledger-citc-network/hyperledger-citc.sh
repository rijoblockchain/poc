#!/bin/bash
sudo apt-get update
sudo apt --yes install software-properties-common
#sudo add-apt-repository --yes ppa:deadsnakes/ppa
#sudo apt install -y python3-pip
#sudo apt install -y build-essential libssl-dev libffi-dev python3-dev
#sudo apt-get --yes install python-tk
#sudo apt-get --yes install iperf
sudo apt install --yes golang
sudo apt-get --yes install screen
sudo apt-get --yes install htop
sudo apt-get install -y jq
sudo apt-get --yes install maven
sudo apt install --yes default-jdk
sudo apt-get install -y sysstat
sudo apt-get install -y curl
sudo apt install docker --yes
sudo apt install docker.io --yes
sudo apt install docker-compose --yes
sudo apt install nodejs --yes
sudo apt install npm --yes
#sudo docker exec -it fablo-rest.org1.example.com  bash

#sudo curl -Lf https://github.com/hyperledger-labs/fablo/releases/download/1.1.0/fablo.sh -o /usr/local/bin/fablo && sudo chmod +x /usr/local/bin/fablo
#sudo fablo init rest dev

# Downs the network and removes fablo-target directory

#sudo fablo prune 

#down and up steps combined. Network state is lost
#sudo fablo reset

#when editing config
#sudo fablo recreate

#validate the network
#sudo fablo validate
