const { NotesService } = require("m3o/notes");

// Specify the note to events
async function subscribeToEvents() {
  let notesService = new NotesService(process.env.MICRO_API_TOKEN);
  let rsp = await notesService.events({
    id: "63c0cdf8-2121-11ec-a881-0242e36f037a",
  });
  console.log(rsp);
}

subscribeToEvents();
