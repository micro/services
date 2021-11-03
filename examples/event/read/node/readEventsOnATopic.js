const { EventService } = require("m3o/event");

// Read stored events
async function readEventsOnAtopic() {
  let eventService = new EventService(process.env.MICRO_API_TOKEN);
  let rsp = await eventService.read({
    topic: "user",
  });
  console.log(rsp);
}

readEventsOnAtopic();
