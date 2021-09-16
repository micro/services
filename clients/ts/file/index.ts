import * as m3o from "@m3o/m3o-node";

export class FileService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Delete a file by project name/path
  delete(request: DeleteRequest): Promise<DeleteResponse> {
    return this.client.call(
      "file",
      "Delete",
      request
    ) as Promise<DeleteResponse>;
  }
  // List files by their project and optionally a path.
  list(request: ListRequest): Promise<ListResponse> {
    return this.client.call("file", "List", request) as Promise<ListResponse>;
  }
  // Read a file by path
  read(request: ReadRequest): Promise<ReadResponse> {
    return this.client.call("file", "Read", request) as Promise<ReadResponse>;
  }
  // Save a file
  save(request: SaveRequest): Promise<SaveResponse> {
    return this.client.call("file", "Save", request) as Promise<SaveResponse>;
  }
}

export interface BatchSaveRequest {
  files?: Record[];
}

export interface BatchSaveResponse {}

export interface DeleteRequest {
  // Path to the file
  path?: string;
  // The project name
  project?: string;
}

export interface DeleteResponse {}

export interface ListRequest {
  // Defaults to '/', ie. lists all files in a project.
  // Supply path to a folder if you want to list
  // files inside that folder
  // eg. '/docs'
  path?: string;
  // Project, required for listing.
  project?: string;
}

export interface ListResponse {
  files?: Record[];
}

export interface ReadRequest {
  // Path to the file
  path?: string;
  // Project name
  project?: string;
}

export interface ReadResponse {
  // Returns the file
  file?: Record;
}

export interface Record {
  // File contents
  content?: string;
  // Time the file was created e.g 2021-05-20T13:37:21Z
  created?: string;
  // Any other associated metadata as a map of key-value pairs
  metadata?: { [key: string]: string };
  // Path to file or folder eg. '/documents/text-files/file.txt'.
  path?: string;
  // A custom project to group files
  // eg. file-of-mywebsite.com
  project?: string;
  // Time the file was updated e.g 2021-05-20T13:37:21Z
  updated?: string;
}

export interface SaveRequest {
  file?: Record;
}

export interface SaveResponse {}
