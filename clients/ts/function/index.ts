import * as m3o from "@m3o/m3o-node";

export class FunctionService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Call a function by name
  call(request: CallRequest): Promise<CallResponse> {
    return this.client.call(
      "function",
      "Call",
      request
    ) as Promise<CallResponse>;
  }
  // Delete a function by name
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
  // List all the deployed functions
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
  // The name of the function
  name?: string;
  // Optional project name
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
  // runtime/language of the function
  // eg: php74,
  // nodejs6, nodejs8, nodejs10, nodejs12, nodejs14, nodejs16
  // dotnet3
  // java11
  // ruby26, ruby27
  // go111, go113, go116
  // python37, python38, python39
  runtime?: string;
  // optional subfolder path
  subfolder?: string;
}

export interface DeployResponse {}

export interface Func {
  // name of handler in source code
  entrypoint?: string;
  // function name
  name?: string;
  // project of function, optional
  // defaults to literal "default"
  // used to namespace functions
  project?: string;
  // git repo address
  repo?: string;
  // runtime/language of the function
  // eg: php74,
  // nodejs6, nodejs8, nodejs10, nodejs12, nodejs14, nodejs16
  // dotnet3
  // java11
  // ruby26, ruby27
  // go111, go113, go116
  // python37, python38, python39
  runtime?: string;
  // subfolder path to entrypoint
  subfolder?: string;
}

export interface ListRequest {
  // optional project name
  project?: string;
}

export interface ListResponse {
  // List of functions deployed
  functions?: Func[];
}
