import * as file from "m3o/file";

// Read a file by path
async function ReadFile() {
  let fileService = new file.FileService(process.env.MICRO_API_TOKEN);
  let rsp = await fileService.read({
    path: "/document/text-files/file.txt",
    project: "examples",
  });
  console.log(rsp);
}

await ReadFile();
