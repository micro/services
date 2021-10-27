const { DbService } = require("m3o/db");

// Create a record in the database. Optionally include an "id" field otherwise it's set automatically.
async function createArecord() {
  let dbService = new DbService(process.env.MICRO_API_TOKEN);
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

createArecord();
