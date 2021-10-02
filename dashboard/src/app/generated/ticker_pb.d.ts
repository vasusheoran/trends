// package: proto
// file: ticker.proto

import * as jspb from "google-protobuf";

export class StockRequest extends jspb.Message {
  getSymbol(): string;
  setSymbol(value: string): void;

  getClose(): number;
  setClose(value: number): void;

  getHigh(): number;
  setHigh(value: number): void;

  getLow(): number;
  setLow(value: number): void;

  getDate(): string;
  setDate(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StockRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StockRequest): StockRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StockRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StockRequest;
  static deserializeBinaryFromReader(message: StockRequest, reader: jspb.BinaryReader): StockRequest;
}

export namespace StockRequest {
  export type AsObject = {
    symbol: string,
    close: number,
    high: number,
    low: number,
    date: string,
  }
}

export class StockResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StockResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StockResponse): StockResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StockResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StockResponse;
  static deserializeBinaryFromReader(message: StockResponse, reader: jspb.BinaryReader): StockResponse;
}

export namespace StockResponse {
  export type AsObject = {
  }
}

export class SummaryRequest extends jspb.Message {
  getSas(): string;
  setSas(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SummaryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SummaryRequest): SummaryRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SummaryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SummaryRequest;
  static deserializeBinaryFromReader(message: SummaryRequest, reader: jspb.BinaryReader): SummaryRequest;
}

export namespace SummaryRequest {
  export type AsObject = {
    sas: string,
  }
}

export class SummaryReply extends jspb.Message {
  getClose(): number;
  setClose(value: number): void;

  getHigh(): number;
  setHigh(value: number): void;

  getLow(): number;
  setLow(value: number): void;

  getAverage(): number;
  setAverage(value: number): void;

  getMinlp3(): number;
  setMinlp3(value: number): void;

  getEma5(): number;
  setEma5(value: number): void;

  getEma20(): number;
  setEma20(value: number): void;

  getRsi(): number;
  setRsi(value: number): void;

  getHl3(): number;
  setHl3(value: number): void;

  getTrend(): number;
  setTrend(value: number): void;

  getBuy(): number;
  setBuy(value: number): void;

  getSupport(): number;
  setSupport(value: number): void;

  getBullish(): number;
  setBullish(value: number): void;

  getBarish(): number;
  setBarish(value: number): void;

  getPreviousbuy(): number;
  setPreviousbuy(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SummaryReply.AsObject;
  static toObject(includeInstance: boolean, msg: SummaryReply): SummaryReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SummaryReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SummaryReply;
  static deserializeBinaryFromReader(message: SummaryReply, reader: jspb.BinaryReader): SummaryReply;
}

export namespace SummaryReply {
  export type AsObject = {
    close: number,
    high: number,
    low: number,
    average: number,
    minlp3: number,
    ema5: number,
    ema20: number,
    rsi: number,
    hl3: number,
    trend: number,
    buy: number,
    support: number,
    bullish: number,
    barish: number,
    previousbuy: number,
  }
}

