# Image

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Image/api](https://m3o.com/Image/api).

Endpoints:

## Upload

Upload an image by either sending a base64 encoded image to this endpoint or a URL.
To resize an image before uploading, see the Resize endpoint.


[https://m3o.com/image/api#Upload](https://m3o.com/image/api#Upload)

```js
const { ImageService } = require('m3o/image');

// Upload an image by either sending a base64 encoded image to this endpoint or a URL.
// To resize an image before uploading, see the Resize endpoint.
async function uploadAbase64imageToMicrosCdn() {
	let imageService = new ImageService(process.env.MICRO_API_TOKEN)
	let rsp = await imageService.upload({
  "base64": "data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAYAAAAeP4ixAAAAx0lEQVR4nOzaMaoDMQyE4ZHj+x82vVdhwQoTkzKQEcwP5r0ihT7sbjUTeAJ4HCegXQJYfOYefOyjDuBiz3yjwJBoCIl6QZOeUjTC1Ix1IxEJXF9+0KWsf2bD4bn37OO/c/wuQ9QyRC1D1DJELUPUMkQtQ9QyRC1D1DJELUPUMkQtQ9QyRC1D1DJELUPUMkQtQ9Sa/NG94Tf3j4WBdaxudMEkn4IM2rZBA0wBrvo7aOcpj2emXvLeVt0IGm0GVXUj91mvAAAA//+V2CZl+4AKXwAAAABJRU5ErkJggg==",
  "name": "cat.jpeg",
  "outputURL": true
})
	console.log(rsp)
}

uploadAbase64imageToMicrosCdn()
```
## Upload

Upload an image by either sending a base64 encoded image to this endpoint or a URL.
To resize an image before uploading, see the Resize endpoint.


[https://m3o.com/image/api#Upload](https://m3o.com/image/api#Upload)

```js
const { ImageService } = require('m3o/image');

// Upload an image by either sending a base64 encoded image to this endpoint or a URL.
// To resize an image before uploading, see the Resize endpoint.
async function uploadAnImageFromAurlToMicrosCdn() {
	let imageService = new ImageService(process.env.MICRO_API_TOKEN)
	let rsp = await imageService.upload({
  "name": "cat.jpeg",
  "url": "somewebsite.com/cat.png"
})
	console.log(rsp)
}

uploadAnImageFromAurlToMicrosCdn()
```
## Resize

Resize an image on the fly without storing it (by sending and receiving a base64 encoded image), or resize and upload depending on parameters.
If one of width or height is 0, the image aspect ratio is preserved.
Optional cropping.


[https://m3o.com/image/api#Resize](https://m3o.com/image/api#Resize)

```js
const { ImageService } = require('m3o/image');

// Resize an image on the fly without storing it (by sending and receiving a base64 encoded image), or resize and upload depending on parameters.
// If one of width or height is 0, the image aspect ratio is preserved.
// Optional cropping.
async function base64toHostedImage() {
	let imageService = new ImageService(process.env.MICRO_API_TOKEN)
	let rsp = await imageService.resize({
  "base64": "data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==",
  "height": 100,
  "name": "cat.png",
  "outputURL": true,
  "width": 100
})
	console.log(rsp)
}

base64toHostedImage()
```
## Resize

Resize an image on the fly without storing it (by sending and receiving a base64 encoded image), or resize and upload depending on parameters.
If one of width or height is 0, the image aspect ratio is preserved.
Optional cropping.


[https://m3o.com/image/api#Resize](https://m3o.com/image/api#Resize)

```js
const { ImageService } = require('m3o/image');

// Resize an image on the fly without storing it (by sending and receiving a base64 encoded image), or resize and upload depending on parameters.
// If one of width or height is 0, the image aspect ratio is preserved.
// Optional cropping.
async function base64toBase64image() {
	let imageService = new ImageService(process.env.MICRO_API_TOKEN)
	let rsp = await imageService.resize({
  "base64": "data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==",
  "height": 100,
  "width": 100
})
	console.log(rsp)
}

base64toBase64image()
```
## Resize

Resize an image on the fly without storing it (by sending and receiving a base64 encoded image), or resize and upload depending on parameters.
If one of width or height is 0, the image aspect ratio is preserved.
Optional cropping.


[https://m3o.com/image/api#Resize](https://m3o.com/image/api#Resize)

```js
const { ImageService } = require('m3o/image');

// Resize an image on the fly without storing it (by sending and receiving a base64 encoded image), or resize and upload depending on parameters.
// If one of width or height is 0, the image aspect ratio is preserved.
// Optional cropping.
async function base64toBase64imageWithCropping() {
	let imageService = new ImageService(process.env.MICRO_API_TOKEN)
	let rsp = await imageService.resize({
  "base64": "data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==",
  "cropOptions": {
    "height": 50,
    "width": 50
  },
  "height": 100,
  "width": 100
})
	console.log(rsp)
}

base64toBase64imageWithCropping()
```
## Convert

Convert an image from one format (jpeg, png etc.) to an other either on the fly (from base64 to base64),
or by uploading the conversion result.


[https://m3o.com/image/api#Convert](https://m3o.com/image/api#Convert)

```js
const { ImageService } = require('m3o/image');

// Convert an image from one format (jpeg, png etc.) to an other either on the fly (from base64 to base64),
// or by uploading the conversion result.
async function convertApngImageToAjpegTakenFromAurlAndSavedToAurlOnMicrosCdn() {
	let imageService = new ImageService(process.env.MICRO_API_TOKEN)
	let rsp = await imageService.convert({
  "name": "cat.jpeg",
  "outputURL": true,
  "url": "somewebsite.com/cat.png"
})
	console.log(rsp)
}

convertApngImageToAjpegTakenFromAurlAndSavedToAurlOnMicrosCdn()
```
