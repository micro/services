const { GeocodingService } = require("m3o/geocoding");

// Reverse lookup an address from gps coordinates
async function reverseGeocodeLocation() {
  let geocodingService = new GeocodingService(process.env.MICRO_API_TOKEN);
  let rsp = await geocodingService.reverse({
    latitude: 51.5123064,
    longitude: -0.1216235,
  });
  console.log(rsp);
}

reverseGeocodeLocation();
