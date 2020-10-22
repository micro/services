import * as jspb from "google-protobuf"

export class Post extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getTitle(): string;
  setTitle(value: string): void;

  getSlug(): string;
  setSlug(value: string): void;

  getContent(): string;
  setContent(value: string): void;

  getCreated(): number;
  setCreated(value: number): void;

  getUpdated(): number;
  setUpdated(value: number): void;

  getAuthor(): string;
  setAuthor(value: string): void;

  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  clearTagsList(): void;
  addTags(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Post.AsObject;
  static toObject(includeInstance: boolean, msg: Post): Post.AsObject;
  static serializeBinaryToWriter(message: Post, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Post;
  static deserializeBinaryFromReader(message: Post, reader: jspb.BinaryReader): Post;
}

export namespace Post {
  export type AsObject = {
    id: string,
    title: string,
    slug: string,
    content: string,
    created: number,
    updated: number,
    author: string,
    tagsList: Array<string>,
  }
}

export class QueryRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getSlug(): string;
  setSlug(value: string): void;

  getTag(): string;
  setTag(value: string): void;

  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: QueryRequest): QueryRequest.AsObject;
  static serializeBinaryToWriter(message: QueryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryRequest;
  static deserializeBinaryFromReader(message: QueryRequest, reader: jspb.BinaryReader): QueryRequest;
}

export namespace QueryRequest {
  export type AsObject = {
    id: string,
    slug: string,
    tag: string,
    offset: number,
    limit: number,
  }
}

export class QueryResponse extends jspb.Message {
  getPostsList(): Array<Post>;
  setPostsList(value: Array<Post>): void;
  clearPostsList(): void;
  addPosts(value?: Post, index?: number): Post;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: QueryResponse): QueryResponse.AsObject;
  static serializeBinaryToWriter(message: QueryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryResponse;
  static deserializeBinaryFromReader(message: QueryResponse, reader: jspb.BinaryReader): QueryResponse;
}

export namespace QueryResponse {
  export type AsObject = {
    postsList: Array<Post.AsObject>,
  }
}

export class SaveRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getTitle(): string;
  setTitle(value: string): void;

  getSlug(): string;
  setSlug(value: string): void;

  getContent(): string;
  setContent(value: string): void;

  getTimestamp(): number;
  setTimestamp(value: number): void;

  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  clearTagsList(): void;
  addTags(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SaveRequest): SaveRequest.AsObject;
  static serializeBinaryToWriter(message: SaveRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveRequest;
  static deserializeBinaryFromReader(message: SaveRequest, reader: jspb.BinaryReader): SaveRequest;
}

export namespace SaveRequest {
  export type AsObject = {
    id: string,
    title: string,
    slug: string,
    content: string,
    timestamp: number,
    tagsList: Array<string>,
  }
}

export class SaveResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SaveResponse): SaveResponse.AsObject;
  static serializeBinaryToWriter(message: SaveResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveResponse;
  static deserializeBinaryFromReader(message: SaveResponse, reader: jspb.BinaryReader): SaveResponse;
}

export namespace SaveResponse {
  export type AsObject = {
    id: string,
  }
}

export class DeleteRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRequest): DeleteRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRequest;
  static deserializeBinaryFromReader(message: DeleteRequest, reader: jspb.BinaryReader): DeleteRequest;
}

export namespace DeleteRequest {
  export type AsObject = {
    id: string,
  }
}

export class DeleteResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteResponse): DeleteResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteResponse;
  static deserializeBinaryFromReader(message: DeleteResponse, reader: jspb.BinaryReader): DeleteResponse;
}

export namespace DeleteResponse {
  export type AsObject = {
  }
}

