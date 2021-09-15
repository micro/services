import * as location from "@m3o/services/location";

// Search for entities in a given radius
async function SearchForLocations() {
  let locationService = new location.LocationService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await locationService.search({
    center: {
      latitude: 51.511061,
      longitude: -0.120022,
    },
    numEntities: 10,
    radius: 100,
    type: "bike",
  });
  console.log(rsp);
}

await SearchForLocations();
