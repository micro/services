import * as db from "@m3o/services/db";

// Create a record in the database. Optionally include an "id" field otherwise it's set automatically.
async function CreateArecord() {
  let dbService = new db.DbService(process.env.MICRO_API_TOKEN);
  let rsp = await dbService.create({
    record: {
      age: 42,
      id: "1",
      isActive: true,
      name: "Jane",
    },
    table: "users",
  });
  console.log(rsp);
}

await CreateArecord();
