import * as grpcWeb from 'grpc-web';

import {
  DeleteRequest,
  DeleteResponse,
  QueryRequest,
  QueryResponse,
  SaveRequest,
  SaveResponse} from './posts_pb';

export class PostsClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  query(
    request: QueryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: QueryResponse) => void
  ): grpcWeb.ClientReadableStream<QueryResponse>;

  save(
    request: SaveRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveResponse) => void
  ): grpcWeb.ClientReadableStream<SaveResponse>;

  delete(
    request: DeleteRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteResponse) => void
  ): grpcWeb.ClientReadableStream<DeleteResponse>;

}

export class PostsPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  query(
    request: QueryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<QueryResponse>;

  save(
    request: SaveRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveResponse>;

  delete(
    request: DeleteRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteResponse>;

}

