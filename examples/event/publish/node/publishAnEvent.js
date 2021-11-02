const { EventService } = require("m3o/event");

// Publish a event to the event stream.
async function publishAnEvent() {
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

publishAnEvent();
