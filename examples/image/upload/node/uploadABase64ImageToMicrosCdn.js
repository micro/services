const { ImageService } = require("m3o/image");

// Upload an image by either sending a base64 encoded image to this endpoint or a URL.
// To resize an image before uploading, see the Resize endpoint.
async function uploadAbase64imageToMicrosCdn() {
  let imageService = new ImageService(process.env.MICRO_API_TOKEN);
  let rsp = await imageService.upload({
    base64:
      "data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAYAAAAeP4ixAAAAx0lEQVR4nOzaMaoDMQyE4ZHj+x82vVdhwQoTkzKQEcwP5r0ihT7sbjUTeAJ4HCegXQJYfOYefOyjDuBiz3yjwJBoCIl6QZOeUjTC1Ix1IxEJXF9+0KWsf2bD4bn37OO/c/wuQ9QyRC1D1DJELUPUMkQtQ9QyRC1D1DJELUPUMkQtQ9QyRC1D1DJELUPUMkQtQ9Sa/NG94Tf3j4WBdaxudMEkn4IM2rZBA0wBrvo7aOcpj2emXvLeVt0IGm0GVXUj91mvAAAA//+V2CZl+4AKXwAAAABJRU5ErkJggg==",
    name: "cat.jpeg",
    outputURL: true,
  });
  console.log(rsp);
}

uploadAbase64imageToMicrosCdn();
