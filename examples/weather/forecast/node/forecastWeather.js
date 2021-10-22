const { WeatherService } = require("m3o/weather");

// Get the weather forecast for the next 1-10 days
async function forecastWeather() {
  let weatherService = new WeatherService(process.env.MICRO_API_TOKEN);
  let rsp = await weatherService.forecast({
    days: 2,
    location: "London",
  });
  console.log(rsp);
}

forecastWeather();
