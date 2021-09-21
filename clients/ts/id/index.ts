import * as m3o from "@m3o/m3o-node";

export class IdService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Generate a unique ID. Defaults to uuid.
  generate(request: GenerateRequest): Promise<GenerateResponse> {
    return this.client.call(
      "id",
      "Generate",
      request
    ) as Promise<GenerateResponse>;
  }
  // List the types of IDs available. No query params needed.
  types(request: TypesRequest): Promise<TypesResponse> {
    return this.client.call("id", "Types", request) as Promise<TypesResponse>;
  }
}

export interface GenerateRequest {
  // type of id e.g uuid, shortid, snowflake (64 bit), bigflake (128 bit)
  type?: string;
}

export interface GenerateResponse {
  // the unique id generated
  id?: string;
  // the type of id generated
  type?: string;
}

export interface TypesRequest {}

export interface TypesResponse {
  types?: string[];
}
