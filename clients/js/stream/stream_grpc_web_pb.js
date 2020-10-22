/**
 * @fileoverview gRPC-Web generated client stub for stream
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.stream = require('./stream_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.stream.RouteGuideClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'binary';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.stream.RouteGuidePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'binary';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.stream.Point,
 *   !proto.stream.Feature>}
 */
const methodDescriptor_RouteGuide_GetFeature = new grpc.web.MethodDescriptor(
  '/stream.RouteGuide/GetFeature',
  grpc.web.MethodType.UNARY,
  proto.stream.Point,
  proto.stream.Feature,
  /**
   * @param {!proto.stream.Point} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.stream.Feature.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.stream.Point,
 *   !proto.stream.Feature>}
 */
const methodInfo_RouteGuide_GetFeature = new grpc.web.AbstractClientBase.MethodInfo(
  proto.stream.Feature,
  /**
   * @param {!proto.stream.Point} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.stream.Feature.deserializeBinary
);


/**
 * @param {!proto.stream.Point} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.stream.Feature)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.stream.Feature>|undefined}
 *     The XHR Node Readable Stream
 */
proto.stream.RouteGuideClient.prototype.getFeature =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/stream.RouteGuide/GetFeature',
      request,
      metadata || {},
      methodDescriptor_RouteGuide_GetFeature,
      callback);
};


/**
 * @param {!proto.stream.Point} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.stream.Feature>}
 *     A native promise that resolves to the response
 */
proto.stream.RouteGuidePromiseClient.prototype.getFeature =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/stream.RouteGuide/GetFeature',
      request,
      metadata || {},
      methodDescriptor_RouteGuide_GetFeature);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.stream.Rectangle,
 *   !proto.stream.Feature>}
 */
const methodDescriptor_RouteGuide_ListFeatures = new grpc.web.MethodDescriptor(
  '/stream.RouteGuide/ListFeatures',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.stream.Rectangle,
  proto.stream.Feature,
  /**
   * @param {!proto.stream.Rectangle} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.stream.Feature.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.stream.Rectangle,
 *   !proto.stream.Feature>}
 */
const methodInfo_RouteGuide_ListFeatures = new grpc.web.AbstractClientBase.MethodInfo(
  proto.stream.Feature,
  /**
   * @param {!proto.stream.Rectangle} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.stream.Feature.deserializeBinary
);


/**
 * @param {!proto.stream.Rectangle} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.stream.Feature>}
 *     The XHR Node Readable Stream
 */
proto.stream.RouteGuideClient.prototype.listFeatures =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/stream.RouteGuide/ListFeatures',
      request,
      metadata || {},
      methodDescriptor_RouteGuide_ListFeatures);
};


/**
 * @param {!proto.stream.Rectangle} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.stream.Feature>}
 *     The XHR Node Readable Stream
 */
proto.stream.RouteGuidePromiseClient.prototype.listFeatures =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/stream.RouteGuide/ListFeatures',
      request,
      metadata || {},
      methodDescriptor_RouteGuide_ListFeatures);
};


module.exports = proto.stream;

