import * as routing from "m3o/routing";

// Retrieve a route as a simple list of gps points along with total distance and estimated duration
async function GpsPointsForAroute() {
  let routingService = new routing.RoutingService(process.env.MICRO_API_TOKEN);
  let rsp = await routingService.route({
    destination: {
      latitude: 52.529407,
      longitude: 13.397634,
    },
    origin: {
      latitude: 52.517037,
      longitude: 13.38886,
    },
  });
  console.log(rsp);
}

await GpsPointsForAroute();
