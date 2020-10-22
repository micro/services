import * as grpcWeb from 'grpc-web';

import {
  HistoryRequest,
  HistoryResponse,
  Message,
  NewRequest,
  NewResponse,
  SendRequest,
  SendResponse} from './chat_pb';

export class ChatClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  new(
    request: NewRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: NewResponse) => void
  ): grpcWeb.ClientReadableStream<NewResponse>;

  history(
    request: HistoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: HistoryResponse) => void
  ): grpcWeb.ClientReadableStream<HistoryResponse>;

  send(
    request: SendRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SendResponse) => void
  ): grpcWeb.ClientReadableStream<SendResponse>;

}

export class ChatPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  new(
    request: NewRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<NewResponse>;

  history(
    request: HistoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<HistoryResponse>;

  send(
    request: SendRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SendResponse>;

}

