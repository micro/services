# File

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/File/api](https://m3o.com/File/api).

Endpoints:

## Save

Save a file


[https://m3o.com/file/api#Save](https://m3o.com/file/api#Save)

```js
const { FileService } = require('m3o/file');

// Save a file
async function saveFile() {
	let fileService = new FileService(process.env.MICRO_API_TOKEN)
	let rsp = await fileService.save({
  "file": {
    "content": "file content example",
    "path": "/document/text-files/file.txt",
    "project": "examples"
  }
})
	console.log(rsp)
}

saveFile()
```
## List

List files by their project and optionally a path.


[https://m3o.com/file/api#List](https://m3o.com/file/api#List)

```js
const { FileService } = require('m3o/file');

// List files by their project and optionally a path.
async function listFiles() {
	let fileService = new FileService(process.env.MICRO_API_TOKEN)
	let rsp = await fileService.list({
  "project": "examples"
})
	console.log(rsp)
}

listFiles()
```
## Delete

Delete a file by project name/path


[https://m3o.com/file/api#Delete](https://m3o.com/file/api#Delete)

```js
const { FileService } = require('m3o/file');

// Delete a file by project name/path
async function deleteFile() {
	let fileService = new FileService(process.env.MICRO_API_TOKEN)
	let rsp = await fileService.delete({
  "path": "/document/text-files/file.txt",
  "project": "examples"
})
	console.log(rsp)
}

deleteFile()
```
## Read

Read a file by path


[https://m3o.com/file/api#Read](https://m3o.com/file/api#Read)

```js
const { FileService } = require('m3o/file');

// Read a file by path
async function readFile() {
	let fileService = new FileService(process.env.MICRO_API_TOKEN)
	let rsp = await fileService.read({
  "path": "/document/text-files/file.txt",
  "project": "examples"
})
	console.log(rsp)
}

readFile()
```
