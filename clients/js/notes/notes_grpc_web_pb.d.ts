import * as grpcWeb from 'grpc-web';

import {
  CreateRequest,
  CreateResponse,
  DeleteRequest,
  DeleteResponse,
  ListRequest,
  ListResponse,
  UpdateRequest,
  UpdateResponse} from './notes_pb';

export class NotesClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  list(
    request: ListRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ListResponse) => void
  ): grpcWeb.ClientReadableStream<ListResponse>;

  create(
    request: CreateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CreateResponse) => void
  ): grpcWeb.ClientReadableStream<CreateResponse>;

  delete(
    request: DeleteRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteResponse) => void
  ): grpcWeb.ClientReadableStream<DeleteResponse>;

  update(
    request: UpdateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UpdateResponse) => void
  ): grpcWeb.ClientReadableStream<UpdateResponse>;

}

export class NotesPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  list(
    request: ListRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ListResponse>;

  create(
    request: CreateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateResponse>;

  delete(
    request: DeleteRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteResponse>;

  update(
    request: UpdateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UpdateResponse>;

}

