// package: proto
// file: ticker.proto

import * as ticker_pb from "./ticker_pb";
import {grpc} from "@improbable-eng/grpc-web";

type TickerUpdateStock = {
  readonly methodName: string;
  readonly service: typeof Ticker;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof ticker_pb.StockRequest;
  readonly responseType: typeof ticker_pb.StockResponse;
};

type TickerGetSummary = {
  readonly methodName: string;
  readonly service: typeof Ticker;
  readonly requestStream: true;
  readonly responseStream: true;
  readonly requestType: typeof ticker_pb.SummaryRequest;
  readonly responseType: typeof ticker_pb.SummaryResponse;
};

export class Ticker {
  static readonly serviceName: string;
  static readonly UpdateStock: TickerUpdateStock;
  static readonly GetSummary: TickerGetSummary;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class TickerClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  updateStock(
    requestMessage: ticker_pb.StockRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: ticker_pb.StockResponse|null) => void
  ): UnaryResponse;
  updateStock(
    requestMessage: ticker_pb.StockRequest,
    callback: (error: ServiceError|null, responseMessage: ticker_pb.StockResponse|null) => void
  ): UnaryResponse;
  getSummary(metadata?: grpc.Metadata): BidirectionalStream<ticker_pb.SummaryRequest, ticker_pb.SummaryResponse>;
}

