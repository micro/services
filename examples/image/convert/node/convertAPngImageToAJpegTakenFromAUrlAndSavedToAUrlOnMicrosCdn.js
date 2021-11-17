const { ImageService } = require("micro-js-client/image");

// Convert an image from one format (jpeg, png etc.) to an other either on the fly (from base64 to base64),
// or by uploading the conversion result.
// To use the file parameter you need to send the request as a multipart/form-data rather than the usual application/json
// with each parameter as a form field.
async function convertApngImageToAjpegTakenFromAurlAndSavedToAurlOnMicrosCdn() {
  let imageService = new ImageService(process.env.MICRO_API_TOKEN);
  let rsp = await imageService.convert({
    name: "cat.jpeg",
    outputURL: true,
    url: "somewebsite.com/cat.png",
  });
  console.log(rsp);
}

convertApngImageToAjpegTakenFromAurlAndSavedToAurlOnMicrosCdn();
