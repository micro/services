# Streams Service

The streams service provides an event stream, designed for sending messages from a server to mutliple
clients connecting via Websockets. The Token RPC should be called to generate a token for each client,
the clients should then subscribe using the Subscribe RPC.
