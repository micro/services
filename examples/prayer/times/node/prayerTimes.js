import * as prayer from "m3o/prayer";

//
async function PrayerTimes() {
  let prayerService = new prayer.PrayerService(process.env.MICRO_API_TOKEN);
  let rsp = await prayerService.times({
    location: "london",
  });
  console.log(rsp);
}

await PrayerTimes();
