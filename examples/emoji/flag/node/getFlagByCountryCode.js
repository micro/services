const { EmojiService } = require("m3o/emoji");

// Get the flag for a country. Requires country code e.g GB for great britain
async function getFlagByCountryCode() {
  let emojiService = new EmojiService(process.env.MICRO_API_TOKEN);
  let rsp = await emojiService.flag({
    alias: "GB",
  });
  console.log(rsp);
}

getFlagByCountryCode();
