FROM node:latest

MAINTAINER Takakazu Abe <tkkz@hatten.ba>

ARG REPOSITORY_URL=""
ARG BRANCH="master"
ARG WORK_PLACE="/opt/tkkz-bot"

RUN apt-get install -y git && rm -rf /var/lib/apt/lists/*

WORKDIR ${WORK_PLACE}

RUN git clone -q -b ${BRANCH} ${REPOSITORY_URL} ${WORK_PLACE}

CMD ["bin/hubot", "-a", "slack"]
