import * as file from "@m3o/services/file";

// Save a file
async function SaveFile() {
  let fileService = new file.FileService(process.env.MICRO_API_TOKEN);
  let rsp = await fileService.save({
    file: {
      content: "file content example",
      path: "/document/text-files/file.txt",
      project: "examples",
    },
  });
  console.log(rsp);
}

await SaveFile();
