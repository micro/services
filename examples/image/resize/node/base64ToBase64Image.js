import * as image from "m3o/image";

// Resize an image on the fly without storing it (by sending and receiving a base64 encoded image), or resize and upload depending on parameters.
// If one of width or height is 0, the image aspect ratio is preserved.
// Optional cropping.
async function Base64toBase64image() {
  let imageService = new image.ImageService(process.env.MICRO_API_TOKEN);
  let rsp = await imageService.resize({
    base64:
      "data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==",
    height: 100,
    width: 100,
  });
  console.log(rsp);
}

await Base64toBase64image();
