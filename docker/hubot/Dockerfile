FROM node:latest

MAINTAINER Takakazu Abe <tkkz@hatten.ba>

ENV REPOSITORY_URL=""
ENV BRANCH="master"
ENV WORK_PLACE="/opt/tkkz-bot"

RUN apt-get install -y git && rm -rf /var/lib/apt/lists/*

WORKDIR ${WORK_PLACE}

RUN git clone -q -b ${BRANCH} ${REPOSITORY_URL} ${WORK_PLACE}

CMD ["bin/hubot", "-a", "slack"]
