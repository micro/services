# Notes

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Notes/api](https://m3o.com/Notes/api).

Endpoints:

## Delete

Delete a note


[https://m3o.com/notes/api#Delete](https://m3o.com/notes/api#Delete)

```js
const { NotesService } = require('m3o/notes');

// Delete a note
async function deleteAnote() {
	let notesService = new NotesService(process.env.MICRO_API_TOKEN)
	let rsp = await notesService.delete({
  "id": "63c0cdf8-2121-11ec-a881-0242e36f037a"
})
	console.log(rsp)
}

deleteAnote()
```
## Create

Create a new note


[https://m3o.com/notes/api#Create](https://m3o.com/notes/api#Create)

```js
const { NotesService } = require('m3o/notes');

// Create a new note
async function createAnote() {
	let notesService = new NotesService(process.env.MICRO_API_TOKEN)
	let rsp = await notesService.create({
  "text": "This is my note",
  "title": "New Note"
})
	console.log(rsp)
}

createAnote()
```
## Read

Read a note


[https://m3o.com/notes/api#Read](https://m3o.com/notes/api#Read)

```js
const { NotesService } = require('m3o/notes');

// Read a note
async function readAnote() {
	let notesService = new NotesService(process.env.MICRO_API_TOKEN)
	let rsp = await notesService.read({
  "id": "63c0cdf8-2121-11ec-a881-0242e36f037a"
})
	console.log(rsp)
}

readAnote()
```
## List

List all the notes


[https://m3o.com/notes/api#List](https://m3o.com/notes/api#List)

```js
const { NotesService } = require('m3o/notes');

// List all the notes
async function listAllNotes() {
	let notesService = new NotesService(process.env.MICRO_API_TOKEN)
	let rsp = await notesService.list({})
	console.log(rsp)
}

listAllNotes()
```
## Update

Update a note


[https://m3o.com/notes/api#Update](https://m3o.com/notes/api#Update)

```js
const { NotesService } = require('m3o/notes');

// Update a note
async function updateAnote() {
	let notesService = new NotesService(process.env.MICRO_API_TOKEN)
	let rsp = await notesService.update({
  "note": {
    "id": "63c0cdf8-2121-11ec-a881-0242e36f037a",
    "text": "Updated note text",
    "title": "Update Note"
  }
})
	console.log(rsp)
}

updateAnote()
```
