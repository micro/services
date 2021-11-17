const { EventService } = require("micro-js-client/event");

// Consume events from a given topic.
async function consumeFromAtopic() {
  let eventService = new EventService(process.env.MICRO_API_TOKEN);
  let rsp = await eventService.consume({
    topic: "user",
  });
  console.log(rsp);
}

consumeFromAtopic();
