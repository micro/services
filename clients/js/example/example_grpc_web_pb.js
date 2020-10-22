/**
 * @fileoverview gRPC-Web generated client stub for srv.test_example
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.srv = {};
proto.srv.test_example = require('./example_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.srv.test_example.ExampleClient =
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
proto.srv.test_example.ExamplePromiseClient =
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
 *   !proto.srv.test_example.Request,
 *   !proto.srv.test_example.Response>}
 */
const methodDescriptor_Example_TestExpiry = new grpc.web.MethodDescriptor(
  '/srv.test_example.Example/TestExpiry',
  grpc.web.MethodType.UNARY,
  proto.srv.test_example.Request,
  proto.srv.test_example.Response,
  /**
   * @param {!proto.srv.test_example.Request} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.srv.test_example.Response.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.srv.test_example.Request,
 *   !proto.srv.test_example.Response>}
 */
const methodInfo_Example_TestExpiry = new grpc.web.AbstractClientBase.MethodInfo(
  proto.srv.test_example.Response,
  /**
   * @param {!proto.srv.test_example.Request} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.srv.test_example.Response.deserializeBinary
);


/**
 * @param {!proto.srv.test_example.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.srv.test_example.Response)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.srv.test_example.Response>|undefined}
 *     The XHR Node Readable Stream
 */
proto.srv.test_example.ExampleClient.prototype.testExpiry =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/srv.test_example.Example/TestExpiry',
      request,
      metadata || {},
      methodDescriptor_Example_TestExpiry,
      callback);
};


/**
 * @param {!proto.srv.test_example.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.srv.test_example.Response>}
 *     A native promise that resolves to the response
 */
proto.srv.test_example.ExamplePromiseClient.prototype.testExpiry =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/srv.test_example.Example/TestExpiry',
      request,
      metadata || {},
      methodDescriptor_Example_TestExpiry);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.srv.test_example.Request,
 *   !proto.srv.test_example.Response>}
 */
const methodDescriptor_Example_TestList = new grpc.web.MethodDescriptor(
  '/srv.test_example.Example/TestList',
  grpc.web.MethodType.UNARY,
  proto.srv.test_example.Request,
  proto.srv.test_example.Response,
  /**
   * @param {!proto.srv.test_example.Request} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.srv.test_example.Response.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.srv.test_example.Request,
 *   !proto.srv.test_example.Response>}
 */
const methodInfo_Example_TestList = new grpc.web.AbstractClientBase.MethodInfo(
  proto.srv.test_example.Response,
  /**
   * @param {!proto.srv.test_example.Request} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.srv.test_example.Response.deserializeBinary
);


/**
 * @param {!proto.srv.test_example.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.srv.test_example.Response)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.srv.test_example.Response>|undefined}
 *     The XHR Node Readable Stream
 */
proto.srv.test_example.ExampleClient.prototype.testList =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/srv.test_example.Example/TestList',
      request,
      metadata || {},
      methodDescriptor_Example_TestList,
      callback);
};


/**
 * @param {!proto.srv.test_example.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.srv.test_example.Response>}
 *     A native promise that resolves to the response
 */
proto.srv.test_example.ExamplePromiseClient.prototype.testList =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/srv.test_example.Example/TestList',
      request,
      metadata || {},
      methodDescriptor_Example_TestList);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.srv.test_example.Request,
 *   !proto.srv.test_example.Response>}
 */
const methodDescriptor_Example_TestListLimit = new grpc.web.MethodDescriptor(
  '/srv.test_example.Example/TestListLimit',
  grpc.web.MethodType.UNARY,
  proto.srv.test_example.Request,
  proto.srv.test_example.Response,
  /**
   * @param {!proto.srv.test_example.Request} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.srv.test_example.Response.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.srv.test_example.Request,
 *   !proto.srv.test_example.Response>}
 */
const methodInfo_Example_TestListLimit = new grpc.web.AbstractClientBase.MethodInfo(
  proto.srv.test_example.Response,
  /**
   * @param {!proto.srv.test_example.Request} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.srv.test_example.Response.deserializeBinary
);


/**
 * @param {!proto.srv.test_example.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.srv.test_example.Response)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.srv.test_example.Response>|undefined}
 *     The XHR Node Readable Stream
 */
proto.srv.test_example.ExampleClient.prototype.testListLimit =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/srv.test_example.Example/TestListLimit',
      request,
      metadata || {},
      methodDescriptor_Example_TestListLimit,
      callback);
};


/**
 * @param {!proto.srv.test_example.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.srv.test_example.Response>}
 *     A native promise that resolves to the response
 */
proto.srv.test_example.ExamplePromiseClient.prototype.testListLimit =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/srv.test_example.Example/TestListLimit',
      request,
      metadata || {},
      methodDescriptor_Example_TestListLimit);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.srv.test_example.Request,
 *   !proto.srv.test_example.Response>}
 */
const methodDescriptor_Example_TestListOffset = new grpc.web.MethodDescriptor(
  '/srv.test_example.Example/TestListOffset',
  grpc.web.MethodType.UNARY,
  proto.srv.test_example.Request,
  proto.srv.test_example.Response,
  /**
   * @param {!proto.srv.test_example.Request} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.srv.test_example.Response.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.srv.test_example.Request,
 *   !proto.srv.test_example.Response>}
 */
const methodInfo_Example_TestListOffset = new grpc.web.AbstractClientBase.MethodInfo(
  proto.srv.test_example.Response,
  /**
   * @param {!proto.srv.test_example.Request} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.srv.test_example.Response.deserializeBinary
);


/**
 * @param {!proto.srv.test_example.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.srv.test_example.Response)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.srv.test_example.Response>|undefined}
 *     The XHR Node Readable Stream
 */
proto.srv.test_example.ExampleClient.prototype.testListOffset =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/srv.test_example.Example/TestListOffset',
      request,
      metadata || {},
      methodDescriptor_Example_TestListOffset,
      callback);
};


/**
 * @param {!proto.srv.test_example.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.srv.test_example.Response>}
 *     A native promise that resolves to the response
 */
proto.srv.test_example.ExamplePromiseClient.prototype.testListOffset =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/srv.test_example.Example/TestListOffset',
      request,
      metadata || {},
      methodDescriptor_Example_TestListOffset);
};


module.exports = proto.srv.test_example;

