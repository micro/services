import * as emoji from "m3o/emoji";

// Find an emoji by its alias e.g :beer:
async function FindEmoji() {
  let emojiService = new emoji.EmojiService(process.env.MICRO_API_TOKEN);
  let rsp = await emojiService.find({
    alias: ":beer:",
  });
  console.log(rsp);
}

await FindEmoji();
