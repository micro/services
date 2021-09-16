import * as emoji from "@m3o/services/emoji";

// Print text and renders the emojis with aliases e.g
// let's grab a :beer: becomes let's grab a 🍺
async function PrintTextIncludingEmoji() {
  let emojiService = new emoji.EmojiService(process.env.MICRO_API_TOKEN);
  let rsp = await emojiService.print({
    text: "let's grab a :beer:",
  });
  console.log(rsp);
}

await PrintTextIncludingEmoji();
