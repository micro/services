/**
 * @fileoverview gRPC-Web generated client stub for posts
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.posts = require('./posts_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.posts.PostsClient =
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
proto.posts.PostsPromiseClient =
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
 *   !proto.posts.QueryRequest,
 *   !proto.posts.QueryResponse>}
 */
const methodDescriptor_Posts_Query = new grpc.web.MethodDescriptor(
  '/posts.Posts/Query',
  grpc.web.MethodType.UNARY,
  proto.posts.QueryRequest,
  proto.posts.QueryResponse,
  /**
   * @param {!proto.posts.QueryRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.posts.QueryResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.posts.QueryRequest,
 *   !proto.posts.QueryResponse>}
 */
const methodInfo_Posts_Query = new grpc.web.AbstractClientBase.MethodInfo(
  proto.posts.QueryResponse,
  /**
   * @param {!proto.posts.QueryRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.posts.QueryResponse.deserializeBinary
);


/**
 * @param {!proto.posts.QueryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.posts.QueryResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.posts.QueryResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.posts.PostsClient.prototype.query =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/posts.Posts/Query',
      request,
      metadata || {},
      methodDescriptor_Posts_Query,
      callback);
};


/**
 * @param {!proto.posts.QueryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.posts.QueryResponse>}
 *     A native promise that resolves to the response
 */
proto.posts.PostsPromiseClient.prototype.query =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/posts.Posts/Query',
      request,
      metadata || {},
      methodDescriptor_Posts_Query);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.posts.SaveRequest,
 *   !proto.posts.SaveResponse>}
 */
const methodDescriptor_Posts_Save = new grpc.web.MethodDescriptor(
  '/posts.Posts/Save',
  grpc.web.MethodType.UNARY,
  proto.posts.SaveRequest,
  proto.posts.SaveResponse,
  /**
   * @param {!proto.posts.SaveRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.posts.SaveResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.posts.SaveRequest,
 *   !proto.posts.SaveResponse>}
 */
const methodInfo_Posts_Save = new grpc.web.AbstractClientBase.MethodInfo(
  proto.posts.SaveResponse,
  /**
   * @param {!proto.posts.SaveRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.posts.SaveResponse.deserializeBinary
);


/**
 * @param {!proto.posts.SaveRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.posts.SaveResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.posts.SaveResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.posts.PostsClient.prototype.save =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/posts.Posts/Save',
      request,
      metadata || {},
      methodDescriptor_Posts_Save,
      callback);
};


/**
 * @param {!proto.posts.SaveRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.posts.SaveResponse>}
 *     A native promise that resolves to the response
 */
proto.posts.PostsPromiseClient.prototype.save =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/posts.Posts/Save',
      request,
      metadata || {},
      methodDescriptor_Posts_Save);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.posts.DeleteRequest,
 *   !proto.posts.DeleteResponse>}
 */
const methodDescriptor_Posts_Delete = new grpc.web.MethodDescriptor(
  '/posts.Posts/Delete',
  grpc.web.MethodType.UNARY,
  proto.posts.DeleteRequest,
  proto.posts.DeleteResponse,
  /**
   * @param {!proto.posts.DeleteRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.posts.DeleteResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.posts.DeleteRequest,
 *   !proto.posts.DeleteResponse>}
 */
const methodInfo_Posts_Delete = new grpc.web.AbstractClientBase.MethodInfo(
  proto.posts.DeleteResponse,
  /**
   * @param {!proto.posts.DeleteRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.posts.DeleteResponse.deserializeBinary
);


/**
 * @param {!proto.posts.DeleteRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.posts.DeleteResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.posts.DeleteResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.posts.PostsClient.prototype.delete =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/posts.Posts/Delete',
      request,
      metadata || {},
      methodDescriptor_Posts_Delete,
      callback);
};


/**
 * @param {!proto.posts.DeleteRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.posts.DeleteResponse>}
 *     A native promise that resolves to the response
 */
proto.posts.PostsPromiseClient.prototype.delete =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/posts.Posts/Delete',
      request,
      metadata || {},
      methodDescriptor_Posts_Delete);
};


module.exports = proto.posts;

