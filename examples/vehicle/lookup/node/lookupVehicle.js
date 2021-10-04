import * as vehicle from "m3o/vehicle";

// Lookup a UK vehicle by it's registration number
async function LookupVehicle() {
  let vehicleService = new vehicle.VehicleService(process.env.MICRO_API_TOKEN);
  let rsp = await vehicleService.lookup({
    registration: "LC60OTA",
  });
  console.log(rsp);
}

await LookupVehicle();
