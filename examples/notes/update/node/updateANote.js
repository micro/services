import * as notes from "m3o/notes";

// Update a note
async function UpdateAnote() {
  let notesService = new notes.NotesService(process.env.MICRO_API_TOKEN);
  let rsp = await notesService.update({
    note: {
      id: "63c0cdf8-2121-11ec-a881-0242e36f037a",
      text: "Updated note text",
      title: "Update Note",
    },
  });
  console.log(rsp);
}

await UpdateAnote();
