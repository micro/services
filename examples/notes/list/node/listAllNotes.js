import * as notes from "m3o/notes";

// List all the notes
async function ListAllNotes() {
  let notesService = new notes.NotesService(process.env.MICRO_API_TOKEN);
  let rsp = await notesService.list({});
  console.log(rsp);
}

await ListAllNotes();
