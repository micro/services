const { TimeService } = require("m3o/time");

// Get the current time
async function returnsCurrentTimeOptionallyWithLocation() {
  let timeService = new TimeService(process.env.MICRO_API_TOKEN);
  let rsp = await timeService.now({});
  console.log(rsp);
}

returnsCurrentTimeOptionallyWithLocation();
