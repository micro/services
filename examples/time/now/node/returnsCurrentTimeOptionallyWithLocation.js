import * as time from "m3o/time";

// Get the current time
async function ReturnsCurrentTimeOptionallyWithLocation() {
  let timeService = new time.TimeService(process.env.MICRO_API_TOKEN);
  let rsp = await timeService.now({});
  console.log(rsp);
}

await ReturnsCurrentTimeOptionallyWithLocation();
