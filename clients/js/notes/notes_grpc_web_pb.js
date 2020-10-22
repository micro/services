/**
 * @fileoverview gRPC-Web generated client stub for notes
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.notes = require('./notes_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.notes.NotesClient =
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
proto.notes.NotesPromiseClient =
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
 *   !proto.notes.ListRequest,
 *   !proto.notes.ListResponse>}
 */
const methodDescriptor_Notes_List = new grpc.web.MethodDescriptor(
  '/notes.Notes/List',
  grpc.web.MethodType.UNARY,
  proto.notes.ListRequest,
  proto.notes.ListResponse,
  /**
   * @param {!proto.notes.ListRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.notes.ListResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.notes.ListRequest,
 *   !proto.notes.ListResponse>}
 */
const methodInfo_Notes_List = new grpc.web.AbstractClientBase.MethodInfo(
  proto.notes.ListResponse,
  /**
   * @param {!proto.notes.ListRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.notes.ListResponse.deserializeBinary
);


/**
 * @param {!proto.notes.ListRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.notes.ListResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.notes.ListResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.notes.NotesClient.prototype.list =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/notes.Notes/List',
      request,
      metadata || {},
      methodDescriptor_Notes_List,
      callback);
};


/**
 * @param {!proto.notes.ListRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.notes.ListResponse>}
 *     A native promise that resolves to the response
 */
proto.notes.NotesPromiseClient.prototype.list =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/notes.Notes/List',
      request,
      metadata || {},
      methodDescriptor_Notes_List);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.notes.CreateRequest,
 *   !proto.notes.CreateResponse>}
 */
const methodDescriptor_Notes_Create = new grpc.web.MethodDescriptor(
  '/notes.Notes/Create',
  grpc.web.MethodType.UNARY,
  proto.notes.CreateRequest,
  proto.notes.CreateResponse,
  /**
   * @param {!proto.notes.CreateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.notes.CreateResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.notes.CreateRequest,
 *   !proto.notes.CreateResponse>}
 */
const methodInfo_Notes_Create = new grpc.web.AbstractClientBase.MethodInfo(
  proto.notes.CreateResponse,
  /**
   * @param {!proto.notes.CreateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.notes.CreateResponse.deserializeBinary
);


/**
 * @param {!proto.notes.CreateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.notes.CreateResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.notes.CreateResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.notes.NotesClient.prototype.create =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/notes.Notes/Create',
      request,
      metadata || {},
      methodDescriptor_Notes_Create,
      callback);
};


/**
 * @param {!proto.notes.CreateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.notes.CreateResponse>}
 *     A native promise that resolves to the response
 */
proto.notes.NotesPromiseClient.prototype.create =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/notes.Notes/Create',
      request,
      metadata || {},
      methodDescriptor_Notes_Create);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.notes.DeleteRequest,
 *   !proto.notes.DeleteResponse>}
 */
const methodDescriptor_Notes_Delete = new grpc.web.MethodDescriptor(
  '/notes.Notes/Delete',
  grpc.web.MethodType.UNARY,
  proto.notes.DeleteRequest,
  proto.notes.DeleteResponse,
  /**
   * @param {!proto.notes.DeleteRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.notes.DeleteResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.notes.DeleteRequest,
 *   !proto.notes.DeleteResponse>}
 */
const methodInfo_Notes_Delete = new grpc.web.AbstractClientBase.MethodInfo(
  proto.notes.DeleteResponse,
  /**
   * @param {!proto.notes.DeleteRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.notes.DeleteResponse.deserializeBinary
);


/**
 * @param {!proto.notes.DeleteRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.notes.DeleteResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.notes.DeleteResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.notes.NotesClient.prototype.delete =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/notes.Notes/Delete',
      request,
      metadata || {},
      methodDescriptor_Notes_Delete,
      callback);
};


/**
 * @param {!proto.notes.DeleteRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.notes.DeleteResponse>}
 *     A native promise that resolves to the response
 */
proto.notes.NotesPromiseClient.prototype.delete =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/notes.Notes/Delete',
      request,
      metadata || {},
      methodDescriptor_Notes_Delete);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.notes.UpdateRequest,
 *   !proto.notes.UpdateResponse>}
 */
const methodDescriptor_Notes_Update = new grpc.web.MethodDescriptor(
  '/notes.Notes/Update',
  grpc.web.MethodType.UNARY,
  proto.notes.UpdateRequest,
  proto.notes.UpdateResponse,
  /**
   * @param {!proto.notes.UpdateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.notes.UpdateResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.notes.UpdateRequest,
 *   !proto.notes.UpdateResponse>}
 */
const methodInfo_Notes_Update = new grpc.web.AbstractClientBase.MethodInfo(
  proto.notes.UpdateResponse,
  /**
   * @param {!proto.notes.UpdateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.notes.UpdateResponse.deserializeBinary
);


/**
 * @param {!proto.notes.UpdateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.notes.UpdateResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.notes.UpdateResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.notes.NotesClient.prototype.update =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/notes.Notes/Update',
      request,
      metadata || {},
      methodDescriptor_Notes_Update,
      callback);
};


/**
 * @param {!proto.notes.UpdateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.notes.UpdateResponse>}
 *     A native promise that resolves to the response
 */
proto.notes.NotesPromiseClient.prototype.update =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/notes.Notes/Update',
      request,
      metadata || {},
      methodDescriptor_Notes_Update);
};


module.exports = proto.notes;

