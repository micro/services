const { StreamService } = require("m3o/stream");

// Send a message to the stream.
async function sendMessage() {
  let streamService = new StreamService(process.env.MICRO_API_TOKEN);
  let rsp = await streamService.sendMessage({
    channel: "general",
    text: "Hey checkout this tweet https://twitter.com/m3oservices/status/1455291054295498752",
  });
  console.log(rsp);
}

sendMessage();
