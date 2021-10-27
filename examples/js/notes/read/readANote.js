const { NotesService } = require("m3o/notes");

// Read a note
async function readAnote() {
  let notesService = new NotesService(process.env.MICRO_API_TOKEN);
  let rsp = await notesService.read({
    id: "63c0cdf8-2121-11ec-a881-0242e36f037a",
  });
  console.log(rsp);
}

readAnote();
