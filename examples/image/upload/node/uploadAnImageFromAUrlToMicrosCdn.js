const { ImageService } = require("m3o/image");

// Upload an image by either sending a base64 encoded image to this endpoint or a URL.
// To resize an image before uploading, see the Resize endpoint.
async function uploadAnImageFromAurlToMicrosCdn() {
  let imageService = new ImageService(process.env.MICRO_API_TOKEN);
  let rsp = await imageService.upload({
    name: "cat.jpeg",
    url: "somewebsite.com/cat.png",
  });
  console.log(rsp);
}

uploadAnImageFromAurlToMicrosCdn();
