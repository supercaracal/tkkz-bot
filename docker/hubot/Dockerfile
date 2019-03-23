FROM node:latest

MAINTAINER Takakazu Abe <tkkz@hatten.ba>

ARG WORK_DIR="/opt/tkkz-bot"

RUN npm install -q -y -g typescript

WORKDIR ${WORK_DIR}

CMD ["bin/hubot", "-a", "slack"]
