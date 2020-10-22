/**
 * @fileoverview gRPC-Web generated client stub for tags
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.tags = require('./tags_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.tags.TagsClient =
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
proto.tags.TagsPromiseClient =
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
 *   !proto.tags.AddRequest,
 *   !proto.tags.AddResponse>}
 */
const methodDescriptor_Tags_Add = new grpc.web.MethodDescriptor(
  '/tags.Tags/Add',
  grpc.web.MethodType.UNARY,
  proto.tags.AddRequest,
  proto.tags.AddResponse,
  /**
   * @param {!proto.tags.AddRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.tags.AddResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.tags.AddRequest,
 *   !proto.tags.AddResponse>}
 */
const methodInfo_Tags_Add = new grpc.web.AbstractClientBase.MethodInfo(
  proto.tags.AddResponse,
  /**
   * @param {!proto.tags.AddRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.tags.AddResponse.deserializeBinary
);


/**
 * @param {!proto.tags.AddRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.tags.AddResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.tags.AddResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.tags.TagsClient.prototype.add =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/tags.Tags/Add',
      request,
      metadata || {},
      methodDescriptor_Tags_Add,
      callback);
};


/**
 * @param {!proto.tags.AddRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.tags.AddResponse>}
 *     A native promise that resolves to the response
 */
proto.tags.TagsPromiseClient.prototype.add =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/tags.Tags/Add',
      request,
      metadata || {},
      methodDescriptor_Tags_Add);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.tags.RemoveRequest,
 *   !proto.tags.RemoveResponse>}
 */
const methodDescriptor_Tags_Remove = new grpc.web.MethodDescriptor(
  '/tags.Tags/Remove',
  grpc.web.MethodType.UNARY,
  proto.tags.RemoveRequest,
  proto.tags.RemoveResponse,
  /**
   * @param {!proto.tags.RemoveRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.tags.RemoveResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.tags.RemoveRequest,
 *   !proto.tags.RemoveResponse>}
 */
const methodInfo_Tags_Remove = new grpc.web.AbstractClientBase.MethodInfo(
  proto.tags.RemoveResponse,
  /**
   * @param {!proto.tags.RemoveRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.tags.RemoveResponse.deserializeBinary
);


/**
 * @param {!proto.tags.RemoveRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.tags.RemoveResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.tags.RemoveResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.tags.TagsClient.prototype.remove =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/tags.Tags/Remove',
      request,
      metadata || {},
      methodDescriptor_Tags_Remove,
      callback);
};


/**
 * @param {!proto.tags.RemoveRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.tags.RemoveResponse>}
 *     A native promise that resolves to the response
 */
proto.tags.TagsPromiseClient.prototype.remove =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/tags.Tags/Remove',
      request,
      metadata || {},
      methodDescriptor_Tags_Remove);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.tags.ListRequest,
 *   !proto.tags.ListResponse>}
 */
const methodDescriptor_Tags_List = new grpc.web.MethodDescriptor(
  '/tags.Tags/List',
  grpc.web.MethodType.UNARY,
  proto.tags.ListRequest,
  proto.tags.ListResponse,
  /**
   * @param {!proto.tags.ListRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.tags.ListResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.tags.ListRequest,
 *   !proto.tags.ListResponse>}
 */
const methodInfo_Tags_List = new grpc.web.AbstractClientBase.MethodInfo(
  proto.tags.ListResponse,
  /**
   * @param {!proto.tags.ListRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.tags.ListResponse.deserializeBinary
);


/**
 * @param {!proto.tags.ListRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.tags.ListResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.tags.ListResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.tags.TagsClient.prototype.list =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/tags.Tags/List',
      request,
      metadata || {},
      methodDescriptor_Tags_List,
      callback);
};


/**
 * @param {!proto.tags.ListRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.tags.ListResponse>}
 *     A native promise that resolves to the response
 */
proto.tags.TagsPromiseClient.prototype.list =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/tags.Tags/List',
      request,
      metadata || {},
      methodDescriptor_Tags_List);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.tags.UpdateRequest,
 *   !proto.tags.UpdateResponse>}
 */
const methodDescriptor_Tags_Update = new grpc.web.MethodDescriptor(
  '/tags.Tags/Update',
  grpc.web.MethodType.UNARY,
  proto.tags.UpdateRequest,
  proto.tags.UpdateResponse,
  /**
   * @param {!proto.tags.UpdateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.tags.UpdateResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.tags.UpdateRequest,
 *   !proto.tags.UpdateResponse>}
 */
const methodInfo_Tags_Update = new grpc.web.AbstractClientBase.MethodInfo(
  proto.tags.UpdateResponse,
  /**
   * @param {!proto.tags.UpdateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.tags.UpdateResponse.deserializeBinary
);


/**
 * @param {!proto.tags.UpdateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.tags.UpdateResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.tags.UpdateResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.tags.TagsClient.prototype.update =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/tags.Tags/Update',
      request,
      metadata || {},
      methodDescriptor_Tags_Update,
      callback);
};


/**
 * @param {!proto.tags.UpdateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.tags.UpdateResponse>}
 *     A native promise that resolves to the response
 */
proto.tags.TagsPromiseClient.prototype.update =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/tags.Tags/Update',
      request,
      metadata || {},
      methodDescriptor_Tags_Update);
};


module.exports = proto.tags;

