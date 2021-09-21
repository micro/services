import * as weather from "m3o/weather";

// Get the current weather report for a location by postcode, city, zip code, ip address
async function GetCurrentWeather() {
  let weatherService = new weather.WeatherService(process.env.MICRO_API_TOKEN);
  let rsp = await weatherService.now({
    location: "london",
  });
  console.log(rsp);
}

await GetCurrentWeather();
