import * as grpcWeb from 'grpc-web';

import {
  Ping,
  Pong,
  Request,
  Response,
  StreamingRequest,
  StreamingResponse} from './idiomatic_pb';

export class IdiomaticClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  call(
    request: Request,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Response) => void
  ): grpcWeb.ClientReadableStream<Response>;

  stream(
    request: StreamingRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<StreamingResponse>;

}

export class IdiomaticPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  call(
    request: Request,
    metadata?: grpcWeb.Metadata
  ): Promise<Response>;

  stream(
    request: StreamingRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<StreamingResponse>;

}

