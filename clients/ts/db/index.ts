import * as m3o from "@m3o/m3o-node";

export class DbService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Count records in a table
  count(request: CountRequest): Promise<CountResponse> {
    return this.client.call("db", "Count", request) as Promise<CountResponse>;
  }
  // Create a record in the database. Optionally include an "id" field otherwise it's set automatically.
  create(request: CreateRequest): Promise<CreateResponse> {
    return this.client.call("db", "Create", request) as Promise<CreateResponse>;
  }
  // Delete a record in the database by id.
  delete(request: DeleteRequest): Promise<DeleteResponse> {
    return this.client.call("db", "Delete", request) as Promise<DeleteResponse>;
  }
  // Read data from a table. Lookup can be by ID or via querying any field in the record.
  read(request: ReadRequest): Promise<ReadResponse> {
    return this.client.call("db", "Read", request) as Promise<ReadResponse>;
  }
  // Truncate the records in a table
  truncate(request: TruncateRequest): Promise<TruncateResponse> {
    return this.client.call(
      "db",
      "Truncate",
      request
    ) as Promise<TruncateResponse>;
  }
  // Update a record in the database. Include an "id" in the record to update.
  update(request: UpdateRequest): Promise<UpdateResponse> {
    return this.client.call("db", "Update", request) as Promise<UpdateResponse>;
  }
}

export interface CountRequest {
  table?: string;
}

export interface CountResponse {
  count?: number;
}

export interface CreateRequest {
  // JSON encoded record or records (can be array or object)
  record?: { [key: string]: any };
  // Optional table name. Defaults to 'default'
  table?: string;
}

export interface CreateResponse {
  // The id of the record (either specified or automatically created)
  id?: string;
}

export interface DeleteRequest {
  // id of the record
  id?: string;
  // Optional table name. Defaults to 'default'
  table?: string;
}

export interface DeleteResponse {}

export interface ReadRequest {
  // Read by id. Equivalent to 'id == "your-id"'
  id?: string;
  // Maximum number of records to return. Default limit is 25.
  // Maximum limit is 1000. Anything higher will return an error.
  limit?: number;
  offset?: number;
  // 'asc' (default), 'desc'
  order?: string;
  // field name to order by
  orderBy?: string;
  // Examples: 'age >= 18', 'age >= 18 and verified == true'
  // Comparison operators: '==', '!=', '<', '>', '<=', '>='
  // Logical operator: 'and'
  // Dot access is supported, eg: 'user.age == 11'
  // Accessing list elements is not supported yet.
  query?: string;
  // Optional table name. Defaults to 'default'
  table?: string;
}

export interface ReadResponse {
  // JSON encoded records
  records?: { [key: string]: any }[];
}

export interface TruncateRequest {
  // Optional table name. Defaults to 'default'
  table?: string;
}

export interface TruncateResponse {
  // The table truncated
  table?: string;
}

export interface UpdateRequest {
  // The id of the record. If not specified it is inferred from the 'id' field of the record
  id?: string;
  // record, JSON object
  record?: { [key: string]: any };
  // Optional table name. Defaults to 'default'
  table?: string;
}

export interface UpdateResponse {}
