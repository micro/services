/**
 * @fileoverview gRPC-Web generated client stub for messages
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.messages = require('./messages_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.messages.MessagesClient =
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
proto.messages.MessagesPromiseClient =
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
 *   !proto.messages.SendRequest,
 *   !proto.messages.SendResponse>}
 */
const methodDescriptor_Messages_Send = new grpc.web.MethodDescriptor(
  '/messages.Messages/Send',
  grpc.web.MethodType.UNARY,
  proto.messages.SendRequest,
  proto.messages.SendResponse,
  /**
   * @param {!proto.messages.SendRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.messages.SendResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.messages.SendRequest,
 *   !proto.messages.SendResponse>}
 */
const methodInfo_Messages_Send = new grpc.web.AbstractClientBase.MethodInfo(
  proto.messages.SendResponse,
  /**
   * @param {!proto.messages.SendRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.messages.SendResponse.deserializeBinary
);


/**
 * @param {!proto.messages.SendRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.messages.SendResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.messages.SendResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.messages.MessagesClient.prototype.send =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/messages.Messages/Send',
      request,
      metadata || {},
      methodDescriptor_Messages_Send,
      callback);
};


/**
 * @param {!proto.messages.SendRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.messages.SendResponse>}
 *     A native promise that resolves to the response
 */
proto.messages.MessagesPromiseClient.prototype.send =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/messages.Messages/Send',
      request,
      metadata || {},
      methodDescriptor_Messages_Send);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.messages.ListRequest,
 *   !proto.messages.ListResponse>}
 */
const methodDescriptor_Messages_List = new grpc.web.MethodDescriptor(
  '/messages.Messages/List',
  grpc.web.MethodType.UNARY,
  proto.messages.ListRequest,
  proto.messages.ListResponse,
  /**
   * @param {!proto.messages.ListRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.messages.ListResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.messages.ListRequest,
 *   !proto.messages.ListResponse>}
 */
const methodInfo_Messages_List = new grpc.web.AbstractClientBase.MethodInfo(
  proto.messages.ListResponse,
  /**
   * @param {!proto.messages.ListRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.messages.ListResponse.deserializeBinary
);


/**
 * @param {!proto.messages.ListRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.messages.ListResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.messages.ListResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.messages.MessagesClient.prototype.list =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/messages.Messages/List',
      request,
      metadata || {},
      methodDescriptor_Messages_List,
      callback);
};


/**
 * @param {!proto.messages.ListRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.messages.ListResponse>}
 *     A native promise that resolves to the response
 */
proto.messages.MessagesPromiseClient.prototype.list =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/messages.Messages/List',
      request,
      metadata || {},
      methodDescriptor_Messages_List);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.messages.ReadRequest,
 *   !proto.messages.ReadResponse>}
 */
const methodDescriptor_Messages_Read = new grpc.web.MethodDescriptor(
  '/messages.Messages/Read',
  grpc.web.MethodType.UNARY,
  proto.messages.ReadRequest,
  proto.messages.ReadResponse,
  /**
   * @param {!proto.messages.ReadRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.messages.ReadResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.messages.ReadRequest,
 *   !proto.messages.ReadResponse>}
 */
const methodInfo_Messages_Read = new grpc.web.AbstractClientBase.MethodInfo(
  proto.messages.ReadResponse,
  /**
   * @param {!proto.messages.ReadRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.messages.ReadResponse.deserializeBinary
);


/**
 * @param {!proto.messages.ReadRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.messages.ReadResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.messages.ReadResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.messages.MessagesClient.prototype.read =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/messages.Messages/Read',
      request,
      metadata || {},
      methodDescriptor_Messages_Read,
      callback);
};


/**
 * @param {!proto.messages.ReadRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.messages.ReadResponse>}
 *     A native promise that resolves to the response
 */
proto.messages.MessagesPromiseClient.prototype.read =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/messages.Messages/Read',
      request,
      metadata || {},
      methodDescriptor_Messages_Read);
};


module.exports = proto.messages;

