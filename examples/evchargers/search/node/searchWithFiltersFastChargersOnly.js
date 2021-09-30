import * as evchargers from "m3o/evchargers";

// Search by giving a coordinate and a max distance, or bounding box and optional filters
async function SearchWithFiltersFastChargersOnly() {
  let evchargersService = new evchargers.EvchargersService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await evchargersService.search({
    distance: 2000,
    levels: ["3"],
    location: {
      latitude: 51.53336351319885,
      longitude: -0.0252,
    },
    max_results: 2,
  });
  console.log(rsp);
}

await SearchWithFiltersFastChargersOnly();
