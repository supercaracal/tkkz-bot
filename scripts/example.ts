import hubot = require("hubot");

module.exports = (robot: hubot.Robot<any>): void => {
  robot.respond(/hello/i, (msg: hubot.Response<hubot.Robot<any>>) => {
    msg.reply("world!");
  });
};
