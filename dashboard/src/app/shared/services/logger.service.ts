import { Injectable } from '@angular/core';
import { LoggingLevel, ConfigLogService } from './config-log.service';

@Injectable({
  providedIn: 'root'
})
export class LoggerService {

  private _level: LoggingLevel = LoggingLevel.Warnings;

  constructor() {
    this._level = LoggingLevel.Info;
  }

  setLoggingLevel(level:LoggingLevel){
    if(level)
      this._level = level
  }

  private log(message: any, level = LoggingLevel.Warnings, ...optionalParams: any[]) {
    if (this.shouldLog(level)) {
      switch (level) {
        case LoggingLevel.Errors:
          console.error(message, optionalParams);
          break;
        case LoggingLevel.Warnings:
          console.warn(message, optionalParams);
          break;
        case LoggingLevel.Info:
          console.info(message, optionalParams);
          break;
        default:
          console.debug(message, optionalParams);
      }
    }
  }

  private shouldLog(level: LoggingLevel) {
    if (this._level === LoggingLevel.None) {
      return false;
    } else if (this._level === LoggingLevel.Errors) {
      return level === LoggingLevel.Errors;
    } else if (this._level === LoggingLevel.Warnings) {
      return level === LoggingLevel.Errors || level === LoggingLevel.Warnings;
    } else if (this._level === LoggingLevel.Info) {
      return level === LoggingLevel.Errors || level === LoggingLevel.Warnings || level === LoggingLevel.Info;
    } else {
      return true;
    }
  }

  logError(message: any, ...optionalParams: any[]) {
    this.log(message, LoggingLevel.Errors, optionalParams);
  }
  
  logWarning(message: any, ...optionalParams: any[]) {
    this.log(message, LoggingLevel.Warnings, optionalParams);
  }
  
  logInfo(message: any, ...optionalParams: any[]) {
    this.log(message, LoggingLevel.Info, optionalParams);
  }
  
  logVerbose(message: any, ...optionalParams: any[]) {
    this.log(message, LoggingLevel.Verbose, optionalParams);
  }
}
