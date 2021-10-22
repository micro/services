const { LocationService } = require("m3o/location");

// Read an entity by its ID
async function getLocationById() {
  let locationService = new LocationService(process.env.MICRO_API_TOKEN);
  let rsp = await locationService.read({
    id: "1",
  });
  console.log(rsp);
}

getLocationById();
