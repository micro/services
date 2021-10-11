import * as m3o from "@m3o/m3o-node";

export class FunctionService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Call a function
  call(request: CallRequest): Promise<CallResponse> {
    return this.client.call(
      "function",
      "Call",
      request
    ) as Promise<CallResponse>;
  }
  //
  delete(request: DeleteRequest): Promise<DeleteResponse> {
    return this.client.call(
      "function",
      "Delete",
      request
    ) as Promise<DeleteResponse>;
  }
  // Deploy a group of functions
  deploy(request: DeployRequest): Promise<DeployResponse> {
    return this.client.call(
      "function",
      "Deploy",
      request
    ) as Promise<DeployResponse>;
  }
  //
  list(request: ListRequest): Promise<ListResponse> {
    return this.client.call(
      "function",
      "List",
      request
    ) as Promise<ListResponse>;
  }
}

export interface CallRequest {
  // Name of the function
  name?: string;
  // Request body that will be passed to the function
  request?: { [key: string]: any };
}

export interface CallResponse {
  // Response body that the function returned
  response?: { [key: string]: any };
}

export interface DeleteRequest {
  name?: string;
  project?: string;
}

export interface DeleteResponse {}

export interface DeployRequest {
  // entry point, ie. handler name in the source code
  // if not provided, defaults to the name parameter
  entrypoint?: string;
  // function name
  name?: string;
  // project is used for namespacing your functions
  // optional. defaults to "default".
  project?: string;
  // github url to repo
  repo?: string;
  // optional subfolder path
  subfolder?: string;
}

export interface DeployResponse {}

export interface Func {
  entrypoint?: string;
  name?: string;
  project?: string;
  repo?: string;
  subfolder?: string;
}

export interface ListRequest {
  // optional
  project?: string;
}

export interface ListResponse {
  functions?: Func[];
}
