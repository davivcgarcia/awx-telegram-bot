# AWX-Telegram-Bot
> Bot used for operational tasks on my personal lab

This is a simple Telegram Bot, written in Golang, that interacts with Ansible AWX (Tower) through RESTful APIs. The idea is to have a bot that can executes Ansible Playbook by simple Telegram commands.

## Build

To make things easier to run and maintain, you should use it as a container. To build a Docker container image, just use the `Dockerfile` included on this repository (Docker 17.05+) and run:

```bash
docker build . -t awx-telegram-bot
docker image prune -f
```

## Execute

The bot requires some environment variables to run. You should run it with:

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