import * as grpcWeb from 'grpc-web';

import {
  Request,
  Response} from './example_pb';

export class ExampleClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  testExpiry(
    request: Request,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Response) => void
  ): grpcWeb.ClientReadableStream<Response>;

  testList(
    request: Request,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Response) => void
  ): grpcWeb.ClientReadableStream<Response>;

  testListLimit(
    request: Request,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Response) => void
  ): grpcWeb.ClientReadableStream<Response>;

  testListOffset(
    request: Request,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Response) => void
  ): grpcWeb.ClientReadableStream<Response>;

}

export class ExamplePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  testExpiry(
    request: Request,
    metadata?: grpcWeb.Metadata
  ): Promise<Response>;

  testList(
    request: Request,
    metadata?: grpcWeb.Metadata
  ): Promise<Response>;

  testListLimit(
    request: Request,
    metadata?: grpcWeb.Metadata
  ): Promise<Response>;

  testListOffset(
    request: Request,
    metadata?: grpcWeb.Metadata
  ): Promise<Response>;

}

