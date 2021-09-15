import * as time from "@m3o/services/time";

// Get the timezone info for a specific location
async function GetTheTimezoneInfoForAspecificLocation() {
  let timeService = new time.TimeService(process.env.MICRO_API_TOKEN);
  let rsp = await timeService.zone({
    location: "London",
  });
  console.log(rsp);
}

await GetTheTimezoneInfoForAspecificLocation();
