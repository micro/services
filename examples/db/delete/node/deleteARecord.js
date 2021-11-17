const { DbService } = require("micro-js-client/db");

// Delete a record in the database by id.
async function deleteArecord() {
  let dbService = new DbService(process.env.MICRO_API_TOKEN);
  let rsp = await dbService.delete({
    id: "1",
    table: "users",
  });
  console.log(rsp);
}

deleteArecord();
