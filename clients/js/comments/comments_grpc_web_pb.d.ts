import * as grpcWeb from 'grpc-web';

import {
  Request,
  Response} from './comments_pb';

export class CommentsClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  save(
    request: Request,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Response) => void
  ): grpcWeb.ClientReadableStream<Response>;

}

export class CommentsPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  save(
    request: Request,
    metadata?: grpcWeb.Metadata
  ): Promise<Response>;

}

