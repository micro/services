const { MqService } = require("micro-js-client/mq");

// Subscribe to messages for a given topic.
async function subscribeToAtopic() {
  let mqService = new MqService(process.env.MICRO_API_TOKEN);
  let rsp = await mqService.subscribe({
    topic: "events",
  });
  console.log(rsp);
}

subscribeToAtopic();
