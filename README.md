[![Go Report Card](https://goreportcard.com/badge/github.com/davivcgarcia/awx-telegram-bot)](https://goreportcard.com/report/github.com/davivcgarcia/awx-telegram-bot) [![Build Status](https://travis-ci.org/davivcgarcia/awx-telegram-bot.svg?branch=master)](https://travis-ci.org/davivcgarcia/awx-telegram-bot) [![Docker Automated build](https://img.shields.io/docker/automated/jrottenberg/ffmpeg.svg)](https://hub.docker.com/r/davivcgarcia/awx-telegram-bot/) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/davivcgarcia/awx-telegram-bot/blob/master/LICENSE)

# AWX-Telegram-Bot

This is a proof-of-concept of a Telegram Bot, written in Golang, that interacts with Ansible AWX (Tower) through RESTful APIs to run Ansible Playbooks over an Google Cloud lab environment. The bot also controls the lab usage by doing check-in/check-out of users based on their Telegram's usernames, using MongoDB as persistance.

## Build

The application can be build and run using Linux containers. To build the container image, use the `Dockerfile` included on this repository and run the command on any host with Docker 17.05+:

```bash
docker build . -t awx-telegram-bot
docker image prune -f
```

## Execute

The bot requires a MongoDB and some environment variables to run. You should run it with:

```bash
docker run -v /root/container/mongo:/data/db:z
           -p 27017:27017
           -itd --rm mongo

docker run -e TELEGRAM_API_TOKEN="<Telegram Bot API Token>" \
           -e TELEGRAM_CHAT_ID="<Telegram Chat ID>" \
           -e AWX_URL="<AWX URL>" \
           -e AWX_USER="<AWX Username>" \
           -e AWX_PASS="<AWX Password>" \
           -e AWX_JOB_TEMPLATE_START="<Job Template for the Start Playbook>" \
           -e AWX_JOB_TEMPLATE_STOP="<Job Template for the Stop Playbook>" \
           -e MONGO_URL="<MongoDB IP>" \
           -e MONGO_DB="awx-telegram-bot" \
           -e MONGO_COLLECTION="registries" \
           -p 8080:8080 \
           -itd --rm awx-telegram-bot
```
