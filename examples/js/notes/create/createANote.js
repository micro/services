const { NotesService } = require("m3o/notes");

// Create a new note
async function createAnote() {
  let notesService = new NotesService(process.env.MICRO_API_TOKEN);
  let rsp = await notesService.create({
    text: "This is my note",
    title: "New Note",
  });
  console.log(rsp);
}

createAnote();
