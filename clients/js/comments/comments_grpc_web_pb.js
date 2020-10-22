/**
 * @fileoverview gRPC-Web generated client stub for comments
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.comments = require('./comments_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.comments.CommentsClient =
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
proto.comments.CommentsPromiseClient =
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
 *   !proto.comments.Request,
 *   !proto.comments.Response>}
 */
const methodDescriptor_Comments_Save = new grpc.web.MethodDescriptor(
  '/comments.Comments/Save',
  grpc.web.MethodType.UNARY,
  proto.comments.Request,
  proto.comments.Response,
  /**
   * @param {!proto.comments.Request} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.comments.Response.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.comments.Request,
 *   !proto.comments.Response>}
 */
const methodInfo_Comments_Save = new grpc.web.AbstractClientBase.MethodInfo(
  proto.comments.Response,
  /**
   * @param {!proto.comments.Request} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.comments.Response.deserializeBinary
);


/**
 * @param {!proto.comments.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.comments.Response)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.comments.Response>|undefined}
 *     The XHR Node Readable Stream
 */
proto.comments.CommentsClient.prototype.save =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/comments.Comments/Save',
      request,
      metadata || {},
      methodDescriptor_Comments_Save,
      callback);
};


/**
 * @param {!proto.comments.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.comments.Response>}
 *     A native promise that resolves to the response
 */
proto.comments.CommentsPromiseClient.prototype.save =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/comments.Comments/Save',
      request,
      metadata || {},
      methodDescriptor_Comments_Save);
};


module.exports = proto.comments;

