const { EmojiService } = require("m3o/emoji");

// Find an emoji by its alias e.g :beer:
async function findEmoji() {
  let emojiService = new EmojiService(process.env.MICRO_API_TOKEN);
  let rsp = await emojiService.find({
    alias: ":beer:",
  });
  console.log(rsp);
}

findEmoji();
