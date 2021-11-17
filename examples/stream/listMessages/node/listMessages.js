const { StreamService } = require("micro-js-client/stream");

// List messages for a given channel
async function listMessages() {
  let streamService = new StreamService(process.env.MICRO_API_TOKEN);
  let rsp = await streamService.listMessages({
    channel: "general",
  });
  console.log(rsp);
}

listMessages();
