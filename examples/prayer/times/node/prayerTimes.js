const { PrayerService } = require("m3o/prayer");

// Get the prayer (salah) times for a location on a given date
async function prayerTimes() {
  let prayerService = new PrayerService(process.env.MICRO_API_TOKEN);
  let rsp = await prayerService.times({
    location: "london",
  });
  console.log(rsp);
}

prayerTimes();
