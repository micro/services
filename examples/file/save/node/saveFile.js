const { FileService } = require("micro-js-client/file");

// Save a file
async function saveFile() {
  let fileService = new FileService(process.env.MICRO_API_TOKEN);
  let rsp = await fileService.save({
    file: {
      content: "file content example",
      path: "/document/text-files/file.txt",
      project: "examples",
    },
  });
  console.log(rsp);
}

saveFile();
