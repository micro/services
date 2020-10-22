import * as grpcWeb from 'grpc-web';

import {
  Feature,
  Point,
  Rectangle,
  RouteNote,
  RouteSummary} from './stream_pb';

export class RouteGuideClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  getFeature(
    request: Point,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Feature) => void
  ): grpcWeb.ClientReadableStream<Feature>;

  listFeatures(
    request: Rectangle,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<Feature>;

}

export class RouteGuidePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  getFeature(
    request: Point,
    metadata?: grpcWeb.Metadata
  ): Promise<Feature>;

  listFeatures(
    request: Rectangle,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<Feature>;

}

