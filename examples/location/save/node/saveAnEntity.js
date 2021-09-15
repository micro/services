import * as location from "@m3o/services/location";

// Save an entity's current position
async function SaveAnEntity() {
  let locationService = new location.LocationService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await locationService.save({
    entity: {
      id: "1",
      location: {
        latitude: 51.511061,
        longitude: -0.120022,
        timestamp: "1622802761",
      },
      type: "bike",
    },
  });
  console.log(rsp);
}

await SaveAnEntity();
