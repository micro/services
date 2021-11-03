const { StreamService } = require("m3o/stream");

// List all the active channels
async function listChannels() {
  let streamService = new StreamService(process.env.MICRO_API_TOKEN);
  let rsp = await streamService.listChannels({});
  console.log(rsp);
}

listChannels();
