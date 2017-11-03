FROM node:latest

MAINTAINER Takakazu Abe <tkkz@hatten.ba>

ENV REPOSITORY_URL=""
ENV BRANCH="master"
ENV WORK_PLACE="/opt/tkkz-bot"

RUN apt-get install -y --no-install-recommend git \
  && rm -rf /var/lib/apt/lists/*

WORKDIR ${WORK_PLACE}

RUN git clone -q -b ${BRANCH} ${REPOSITORY_URL} ${WORK_PLACE}
RUN npm install -q -y

CMD ["bin/hubot"]
