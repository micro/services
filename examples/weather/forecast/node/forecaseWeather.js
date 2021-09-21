import * as weather from "m3o/weather";

// Get the weather forecast for the next 1-10 days
async function ForecaseWeather() {
  let weatherService = new weather.WeatherService(process.env.MICRO_API_TOKEN);
  let rsp = await weatherService.forecast({
    days: 2,
    location: "London",
  });
  console.log(rsp);
}

await ForecaseWeather();
