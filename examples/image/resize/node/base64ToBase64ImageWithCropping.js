const { ImageService } = require("m3o/image");

// Resize an image on the fly without storing it (by sending and receiving a base64 encoded image), or resize and upload depending on parameters.
// If one of width or height is 0, the image aspect ratio is preserved.
// Optional cropping.
async function base64toBase64imageWithCropping() {
  let imageService = new ImageService(process.env.MICRO_API_TOKEN);
  let rsp = await imageService.resize({
    base64:
      "data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==",
    cropOptions: {
      height: 50,
      width: 50,
    },
    height: 100,
    width: 100,
  });
  console.log(rsp);
}

base64toBase64imageWithCropping();
