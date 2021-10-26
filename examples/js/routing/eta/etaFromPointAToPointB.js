const { RoutingService } = require("m3o/routing");

// Get the eta for a route from origin to destination. The eta is an estimated time based on car routes
async function etaFromPointAtoPointB() {
  let routingService = new RoutingService(process.env.MICRO_API_TOKEN);
  let rsp = await routingService.eta({
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

etaFromPointAtoPointB();
