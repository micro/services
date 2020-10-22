/**
 * @fileoverview gRPC-Web generated client stub for helloworld
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.helloworld = require('./helloworld_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.helloworld.HelloworldClient =
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
proto.helloworld.HelloworldPromiseClient =
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
 *   !proto.helloworld.Request,
 *   !proto.helloworld.Response>}
 */
const methodDescriptor_Helloworld_Call = new grpc.web.MethodDescriptor(
  '/helloworld.Helloworld/Call',
  grpc.web.MethodType.UNARY,
  proto.helloworld.Request,
  proto.helloworld.Response,
  /**
   * @param {!proto.helloworld.Request} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.helloworld.Response.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.helloworld.Request,
 *   !proto.helloworld.Response>}
 */
const methodInfo_Helloworld_Call = new grpc.web.AbstractClientBase.MethodInfo(
  proto.helloworld.Response,
  /**
   * @param {!proto.helloworld.Request} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.helloworld.Response.deserializeBinary
);


/**
 * @param {!proto.helloworld.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.helloworld.Response)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.helloworld.Response>|undefined}
 *     The XHR Node Readable Stream
 */
proto.helloworld.HelloworldClient.prototype.call =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/helloworld.Helloworld/Call',
      request,
      metadata || {},
      methodDescriptor_Helloworld_Call,
      callback);
};


/**
 * @param {!proto.helloworld.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.helloworld.Response>}
 *     A native promise that resolves to the response
 */
proto.helloworld.HelloworldPromiseClient.prototype.call =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/helloworld.Helloworld/Call',
      request,
      metadata || {},
      methodDescriptor_Helloworld_Call);
};


module.exports = proto.helloworld;

