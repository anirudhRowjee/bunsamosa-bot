#!/bin/sh
# ACM PESUECC HackNight BunSamosa Bot eploy Shell Script -> https://tabvn.medium.com/deploy-golang-application-on-digital-ocean-server-ubuntu-16-04-b7bf5340ccd9

## One-Time Production Environment Setup
#sudo apt update
#sudo apt install nginx
#sudo ufw allow 'Nginx HTTP'
#sudo apt update
#sudo apt install software-properties-common
#sudo add-apt-repository universe
## Allow ourselves to provision HTTPS certificates
#sudo add-apt-repository ppa:certbot/certbot
#sudo apt update
#sudo apt install certbot python3-certbot-nginx
#sudo certbot --nginx

GOOS=linux GOARCH=amd64 go build