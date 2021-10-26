const { GeocodingService } = require("m3o/geocoding");

// Lookup returns a geocoded address including normalized address and gps coordinates. All fields are optional, provide more to get more accurate results
async function geocodeAnAddress() {
  let geocodingService = new GeocodingService(process.env.MICRO_API_TOKEN);
  let rsp = await geocodingService.lookup({
    address: "10 russell st",
    city: "london",
    country: "uk",
    postcode: "wc2b",
  });
  console.log(rsp);
}

geocodeAnAddress();
