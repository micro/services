const { StreamService } = require("m3o/stream");

// Publish a message to the stream. Specify a topic to group messages for a specific topic.
async function publishAmessage() {
  let streamService = new StreamService(process.env.MICRO_API_TOKEN);
  let rsp = await streamService.publish({
    message: {
      id: "1",
      type: "signup",
      user: "john",
    },
    topic: "events",
  });
  console.log(rsp);
}

publishAmessage();
