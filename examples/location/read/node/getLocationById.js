import * as location from "@m3o/services/location";

// Read an entity by its ID
async function GetLocationById() {
  let locationService = new location.LocationService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await locationService.read({
    id: "1",
  });
  console.log(rsp);
}

await GetLocationById();
