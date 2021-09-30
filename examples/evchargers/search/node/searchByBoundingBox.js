import * as evchargers from "m3o/evchargers";

// Search by giving a coordinate and a max distance, or bounding box and optional filters
async function SearchByBoundingBox() {
  let evchargersService = new evchargers.EvchargersService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await evchargersService.search({
    box: {
      bottom_left: {
        latitude: 51.52627543859447,
        longitude: -0.03635349400295168,
      },
      top_right: {
        latitude: 51.56717121807993,
        longitude: -0.002293530559768285,
      },
    },
    max_results: 2,
  });
  console.log(rsp);
}

await SearchByBoundingBox();
