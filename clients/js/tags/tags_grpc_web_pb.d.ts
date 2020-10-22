import * as grpcWeb from 'grpc-web';

import {
  AddRequest,
  AddResponse,
  ListRequest,
  ListResponse,
  RemoveRequest,
  RemoveResponse,
  UpdateRequest,
  UpdateResponse} from './tags_pb';

export class TagsClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  add(
    request: AddRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AddResponse) => void
  ): grpcWeb.ClientReadableStream<AddResponse>;

  remove(
    request: RemoveRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RemoveResponse) => void
  ): grpcWeb.ClientReadableStream<RemoveResponse>;

  list(
    request: ListRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ListResponse) => void
  ): grpcWeb.ClientReadableStream<ListResponse>;

  update(
    request: UpdateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UpdateResponse) => void
  ): grpcWeb.ClientReadableStream<UpdateResponse>;

}

export class TagsPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  add(
    request: AddRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<AddResponse>;

  remove(
    request: RemoveRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<RemoveResponse>;

  list(
    request: ListRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ListResponse>;

  update(
    request: UpdateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UpdateResponse>;

}

