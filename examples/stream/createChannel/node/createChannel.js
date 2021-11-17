const { StreamService } = require("micro-js-client/stream");

// Create a channel with a given name and description. Channels are created automatically but
// this allows you to specify a description that's persisted for the lifetime of the channel.
async function createChannel() {
  let streamService = new StreamService(process.env.MICRO_API_TOKEN);
  let rsp = await streamService.createChannel({
    description: "The channel for all things",
    name: "general",
  });
  console.log(rsp);
}

createChannel();
