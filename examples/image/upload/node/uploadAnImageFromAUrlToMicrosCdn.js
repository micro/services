const { ImageService } = require("micro-js-client/image");

// Upload an image by either sending a base64 encoded image to this endpoint or a URL.
// To resize an image before uploading, see the Resize endpoint.
// To use the file parameter you need to send the request as a multipart/form-data rather than the usual application/json
// with each parameter as a form field.
async function uploadAnImageFromAurlToMicrosCdn() {
  let imageService = new ImageService(process.env.MICRO_API_TOKEN);
  let rsp = await imageService.upload({
    name: "cat.jpeg",
    url: "somewebsite.com/cat.png",
  });
  console.log(rsp);
}

uploadAnImageFromAurlToMicrosCdn();
