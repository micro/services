import * as file from "m3o/file";

// List files by their project and optionally a path.
async function ListFiles() {
  let fileService = new file.FileService(process.env.MICRO_API_TOKEN);
  let rsp = await fileService.list({
    project: "examples",
  });
  console.log(rsp);
}

await ListFiles();
