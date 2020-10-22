/**
 * @fileoverview gRPC-Web generated client stub for idiomatic
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.idiomatic = require('./idiomatic_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.idiomatic.IdiomaticClient =
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
proto.idiomatic.IdiomaticPromiseClient =
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
 *   !proto.idiomatic.Request,
 *   !proto.idiomatic.Response>}
 */
const methodDescriptor_Idiomatic_Call = new grpc.web.MethodDescriptor(
  '/idiomatic.Idiomatic/Call',
  grpc.web.MethodType.UNARY,
  proto.idiomatic.Request,
  proto.idiomatic.Response,
  /**
   * @param {!proto.idiomatic.Request} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.idiomatic.Response.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.idiomatic.Request,
 *   !proto.idiomatic.Response>}
 */
const methodInfo_Idiomatic_Call = new grpc.web.AbstractClientBase.MethodInfo(
  proto.idiomatic.Response,
  /**
   * @param {!proto.idiomatic.Request} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.idiomatic.Response.deserializeBinary
);


/**
 * @param {!proto.idiomatic.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.idiomatic.Response)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.idiomatic.Response>|undefined}
 *     The XHR Node Readable Stream
 */
proto.idiomatic.IdiomaticClient.prototype.call =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/idiomatic.Idiomatic/Call',
      request,
      metadata || {},
      methodDescriptor_Idiomatic_Call,
      callback);
};


/**
 * @param {!proto.idiomatic.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.idiomatic.Response>}
 *     A native promise that resolves to the response
 */
proto.idiomatic.IdiomaticPromiseClient.prototype.call =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/idiomatic.Idiomatic/Call',
      request,
      metadata || {},
      methodDescriptor_Idiomatic_Call);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.idiomatic.StreamingRequest,
 *   !proto.idiomatic.StreamingResponse>}
 */
const methodDescriptor_Idiomatic_Stream = new grpc.web.MethodDescriptor(
  '/idiomatic.Idiomatic/Stream',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.idiomatic.StreamingRequest,
  proto.idiomatic.StreamingResponse,
  /**
   * @param {!proto.idiomatic.StreamingRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.idiomatic.StreamingResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.idiomatic.StreamingRequest,
 *   !proto.idiomatic.StreamingResponse>}
 */
const methodInfo_Idiomatic_Stream = new grpc.web.AbstractClientBase.MethodInfo(
  proto.idiomatic.StreamingResponse,
  /**
   * @param {!proto.idiomatic.StreamingRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.idiomatic.StreamingResponse.deserializeBinary
);


/**
 * @param {!proto.idiomatic.StreamingRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.idiomatic.StreamingResponse>}
 *     The XHR Node Readable Stream
 */
proto.idiomatic.IdiomaticClient.prototype.stream =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/idiomatic.Idiomatic/Stream',
      request,
      metadata || {},
      methodDescriptor_Idiomatic_Stream);
};


/**
 * @param {!proto.idiomatic.StreamingRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.idiomatic.StreamingResponse>}
 *     The XHR Node Readable Stream
 */
proto.idiomatic.IdiomaticPromiseClient.prototype.stream =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/idiomatic.Idiomatic/Stream',
      request,
      metadata || {},
      methodDescriptor_Idiomatic_Stream);
};


module.exports = proto.idiomatic;

