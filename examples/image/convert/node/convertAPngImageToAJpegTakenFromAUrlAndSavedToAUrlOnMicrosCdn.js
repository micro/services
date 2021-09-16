import * as image from "@m3o/services/image";

// Convert an image from one format (jpeg, png etc.) to an other either on the fly (from base64 to base64),
// or by uploading the conversion result.
async function ConvertApngImageToAjpegTakenFromAurlAndSavedToAurlOnMicrosCdn() {
  let imageService = new image.ImageService(process.env.MICRO_API_TOKEN);
  let rsp = await imageService.convert({
    name: "cat.jpeg",
    outputURL: true,
    url: "somewebsite.com/cat.png",
  });
  console.log(rsp);
}

await ConvertApngImageToAjpegTakenFromAurlAndSavedToAurlOnMicrosCdn();
