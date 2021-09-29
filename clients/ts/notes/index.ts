import * as m3o from "@m3o/m3o-node";

export class NotesService {
  private client: m3o.Client;

  constructor(token: string) {
    this.client = new m3o.Client({ token: token });
  }
  // Create a new note
  create(request: CreateRequest): Promise<CreateResponse> {
    return this.client.call(
      "notes",
      "Create",
      request
    ) as Promise<CreateResponse>;
  }
  // Delete a note
  delete(request: DeleteRequest): Promise<DeleteResponse> {
    return this.client.call(
      "notes",
      "Delete",
      request
    ) as Promise<DeleteResponse>;
  }
  // List all the notes
  list(request: ListRequest): Promise<ListResponse> {
    return this.client.call("notes", "List", request) as Promise<ListResponse>;
  }
  // Read a note
  read(request: ReadRequest): Promise<ReadResponse> {
    return this.client.call("notes", "Read", request) as Promise<ReadResponse>;
  }
  // Specify the note to events
  subscribe(request: SubscribeRequest): Promise<SubscribeResponse> {
    return this.client.call(
      "notes",
      "Subscribe",
      request
    ) as Promise<SubscribeResponse>;
  }
  // Update a note
  update(request: UpdateRequest): Promise<UpdateResponse> {
    return this.client.call(
      "notes",
      "Update",
      request
    ) as Promise<UpdateResponse>;
  }
}

export interface CreateRequest {
  // note text
  text?: string;
  // note title
  title?: string;
}

export interface CreateResponse {
  // The created note
  note?: { [key: string]: any };
}

export interface DeleteRequest {
  // specify the id of the note
  id?: string;
}

export interface DeleteResponse {
  note?: { [key: string]: any };
}

export interface ListRequest {}

export interface ListResponse {
  // the list of notes
  notes?: Note[];
}

export interface Note {
  // time at which the note was created
  created?: string;
  // unique id for the note, generated if not specified
  id?: string;
  // text within the note
  text?: string;
  // title of the note
  title?: string;
  // time at which the note was updated
  updated?: string;
}

export interface ReadRequest {
  // the note id
  id?: string;
}

export interface ReadResponse {
  // The note
  note?: { [key: string]: any };
}

export interface SubscribeRequest {
  // optionally specify a note id
  id?: string;
}

export interface SubscribeResponse {
  // the event which occured; created, deleted, updated
  event?: string;
  // the note which the operation occured on
  note?: { [key: string]: any };
}

export interface UpdateRequest {
  note?: { [key: string]: any };
}

export interface UpdateResponse {
  note?: { [key: string]: any };
}
