import * as routing from "@m3o/services/routing";

// Turn by turn directions from a start point to an end point including maneuvers and bearings
async function TurnByTurnDirections() {
  let routingService = new routing.RoutingService(process.env.MICRO_API_TOKEN);
  let rsp = await routingService.directions({
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

await TurnByTurnDirections();
