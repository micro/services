import * as grpcWeb from 'grpc-web';

import {
  Request,
  Response} from './helloworld_pb';

export class HelloworldClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  call(
    request: Request,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Response) => void
  ): grpcWeb.ClientReadableStream<Response>;

}

export class HelloworldPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  call(
    request: Request,
    metadata?: grpcWeb.Metadata
  ): Promise<Response>;

}

