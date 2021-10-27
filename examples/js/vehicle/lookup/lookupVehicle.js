const { VehicleService } = require("m3o/vehicle");

// Lookup a UK vehicle by it's registration number
async function lookupVehicle() {
  let vehicleService = new VehicleService(process.env.MICRO_API_TOKEN);
  let rsp = await vehicleService.lookup({
    registration: "LC60OTA",
  });
  console.log(rsp);
}

lookupVehicle();
