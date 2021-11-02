const { EventService } = require("m3o/event");

// Subscribe to messages for a given topic.
async function subscribeToAtopic() {
  let eventService = new EventService(process.env.MICRO_API_TOKEN);
  let rsp = await eventService.subscribe({
    topic: "user",
  });
  console.log(rsp);
}

subscribeToAtopic();
