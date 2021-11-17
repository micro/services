const { ImageService } = require("micro-js-client/image");

// Delete an image previously uploaded.
async function deleteAnUploadedImage() {
  let imageService = new ImageService(process.env.MICRO_API_TOKEN);
  let rsp = await imageService.delete({
    url: "https://cdn.m3ocontent.com/micro/images/micro/41e23b39-48dd-42b6-9738-79a313414bb8/cat.png",
  });
  console.log(rsp);
}

deleteAnUploadedImage();
