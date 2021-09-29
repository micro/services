import * as notes from "m3o/notes";

// Create a new note
async function CreateAnote() {
  let notesService = new notes.NotesService(process.env.MICRO_API_TOKEN);
  let rsp = await notesService.create({
    text: "This is my note",
    title: "New Note",
  });
  console.log(rsp);
}

await CreateAnote();
