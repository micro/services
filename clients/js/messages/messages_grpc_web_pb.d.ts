import * as grpcWeb from 'grpc-web';

import {
  ListRequest,
  ListResponse,
  ReadRequest,
  ReadResponse,
  SendRequest,
  SendResponse} from './messages_pb';

export class MessagesClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  send(
    request: SendRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SendResponse) => void
  ): grpcWeb.ClientReadableStream<SendResponse>;

  list(
    request: ListRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ListResponse) => void
  ): grpcWeb.ClientReadableStream<ListResponse>;

  read(
    request: ReadRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ReadResponse) => void
  ): grpcWeb.ClientReadableStream<ReadResponse>;

}

export class MessagesPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  send(
    request: SendRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SendResponse>;

  list(
    request: ListRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ListResponse>;

  read(
    request: ReadRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ReadResponse>;

}

