// package: proto
// file: ticker.proto

var ticker_pb = require("./ticker_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var Ticker = (function () {
  function Ticker() {}
  Ticker.serviceName = "proto.Ticker";
  return Ticker;
}());

Ticker.UpdateStock = {
  methodName: "UpdateStock",
  service: Ticker,
  requestStream: false,
  responseStream: false,
  requestType: ticker_pb.StockRequest,
  responseType: ticker_pb.StockResponse
};

Ticker.GetSummary = {
  methodName: "GetSummary",
  service: Ticker,
  requestStream: false,
  responseStream: true,
  requestType: ticker_pb.SummaryRequest,
  responseType: ticker_pb.SummaryReply
};

exports.Ticker = Ticker;

function TickerClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

TickerClient.prototype.updateStock = function updateStock(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(Ticker.UpdateStock, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

TickerClient.prototype.getSummary = function getSummary(requestMessage, metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.invoke(Ticker.GetSummary, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onMessage: function (responseMessage) {
      listeners.data.forEach(function (handler) {
        handler(responseMessage);
      });
    },
    onEnd: function (status, statusMessage, trailers) {
      listeners.status.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners.end.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners = null;
    }
  });
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

exports.TickerClient = TickerClient;

