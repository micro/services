const { StreamService } = require("m3o/stream");

// Subscribe to messages for a given topic.
async function subscribeToAtopic() {
  let streamService = new StreamService(process.env.MICRO_API_TOKEN);
  let rsp = await streamService.subscribe({
    topic: "events",
  });
  console.log(rsp);
}

subscribeToAtopic();
