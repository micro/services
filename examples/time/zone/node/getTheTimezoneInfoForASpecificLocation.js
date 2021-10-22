const { TimeService } = require("m3o/time");

// Get the timezone info for a specific location
async function getTheTimezoneInfoForAspecificLocation() {
  let timeService = new TimeService(process.env.MICRO_API_TOKEN);
  let rsp = await timeService.zone({
    location: "London",
  });
  console.log(rsp);
}

getTheTimezoneInfoForAspecificLocation();
