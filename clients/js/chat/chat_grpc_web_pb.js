/**
 * @fileoverview gRPC-Web generated client stub for chat
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.chat = require('./chat_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.chat.ChatClient =
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
proto.chat.ChatPromiseClient =
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
 *   !proto.chat.NewRequest,
 *   !proto.chat.NewResponse>}
 */
const methodDescriptor_Chat_New = new grpc.web.MethodDescriptor(
  '/chat.Chat/New',
  grpc.web.MethodType.UNARY,
  proto.chat.NewRequest,
  proto.chat.NewResponse,
  /**
   * @param {!proto.chat.NewRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.chat.NewResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.chat.NewRequest,
 *   !proto.chat.NewResponse>}
 */
const methodInfo_Chat_New = new grpc.web.AbstractClientBase.MethodInfo(
  proto.chat.NewResponse,
  /**
   * @param {!proto.chat.NewRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.chat.NewResponse.deserializeBinary
);


/**
 * @param {!proto.chat.NewRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.chat.NewResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.chat.NewResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.chat.ChatClient.prototype.new =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/chat.Chat/New',
      request,
      metadata || {},
      methodDescriptor_Chat_New,
      callback);
};


/**
 * @param {!proto.chat.NewRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.chat.NewResponse>}
 *     A native promise that resolves to the response
 */
proto.chat.ChatPromiseClient.prototype.new =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/chat.Chat/New',
      request,
      metadata || {},
      methodDescriptor_Chat_New);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.chat.HistoryRequest,
 *   !proto.chat.HistoryResponse>}
 */
const methodDescriptor_Chat_History = new grpc.web.MethodDescriptor(
  '/chat.Chat/History',
  grpc.web.MethodType.UNARY,
  proto.chat.HistoryRequest,
  proto.chat.HistoryResponse,
  /**
   * @param {!proto.chat.HistoryRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.chat.HistoryResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.chat.HistoryRequest,
 *   !proto.chat.HistoryResponse>}
 */
const methodInfo_Chat_History = new grpc.web.AbstractClientBase.MethodInfo(
  proto.chat.HistoryResponse,
  /**
   * @param {!proto.chat.HistoryRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.chat.HistoryResponse.deserializeBinary
);


/**
 * @param {!proto.chat.HistoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.chat.HistoryResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.chat.HistoryResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.chat.ChatClient.prototype.history =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/chat.Chat/History',
      request,
      metadata || {},
      methodDescriptor_Chat_History,
      callback);
};


/**
 * @param {!proto.chat.HistoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.chat.HistoryResponse>}
 *     A native promise that resolves to the response
 */
proto.chat.ChatPromiseClient.prototype.history =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/chat.Chat/History',
      request,
      metadata || {},
      methodDescriptor_Chat_History);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.chat.SendRequest,
 *   !proto.chat.SendResponse>}
 */
const methodDescriptor_Chat_Send = new grpc.web.MethodDescriptor(
  '/chat.Chat/Send',
  grpc.web.MethodType.UNARY,
  proto.chat.SendRequest,
  proto.chat.SendResponse,
  /**
   * @param {!proto.chat.SendRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.chat.SendResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.chat.SendRequest,
 *   !proto.chat.SendResponse>}
 */
const methodInfo_Chat_Send = new grpc.web.AbstractClientBase.MethodInfo(
  proto.chat.SendResponse,
  /**
   * @param {!proto.chat.SendRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.chat.SendResponse.deserializeBinary
);


/**
 * @param {!proto.chat.SendRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.chat.SendResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.chat.SendResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.chat.ChatClient.prototype.send =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/chat.Chat/Send',
      request,
      metadata || {},
      methodDescriptor_Chat_Send,
      callback);
};


/**
 * @param {!proto.chat.SendRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.chat.SendResponse>}
 *     A native promise that resolves to the response
 */
proto.chat.ChatPromiseClient.prototype.send =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/chat.Chat/Send',
      request,
      metadata || {},
      methodDescriptor_Chat_Send);
};


module.exports = proto.chat;

