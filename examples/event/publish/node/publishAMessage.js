const { EventService } = require("m3o/event");

// Publish a message to the event. Specify a topic to group messages for a specific topic.
async function publishAmessage() {
  let eventService = new EventService(process.env.MICRO_API_TOKEN);
  let rsp = await eventService.publish({
    message: {
      id: "1",
      type: "signup",
      user: "john",
    },
    topic: "user",
  });
  console.log(rsp);
}

publishAmessage();
