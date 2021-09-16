import * as db from "m3o/db";

// Delete a record in the database by id.
async function DeleteArecord() {
  let dbService = new db.DbService(process.env.MICRO_API_TOKEN);
  let rsp = await dbService.delete({
    id: "1",
    table: "users",
  });
  console.log(rsp);
}

await DeleteArecord();
