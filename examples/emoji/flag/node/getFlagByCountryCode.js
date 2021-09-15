import * as emoji from "@m3o/services/emoji";

// Get the flag for a country. Requires country code e.g GB for great britain
async function GetFlagByCountryCode() {
  let emojiService = new emoji.EmojiService(process.env.MICRO_API_TOKEN);
  let rsp = await emojiService.flag({
    alias: "GB",
  });
  console.log(rsp);
}

await GetFlagByCountryCode();
