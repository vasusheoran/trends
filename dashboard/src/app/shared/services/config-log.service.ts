import { Injectable } from '@angular/core';

export class LoggingLevel {
  public static None = 'None';
  public static Verbose = 'Verbose';
  public static Info = 'Info';
  public static Warnings = 'Warnings';
  public static Errors = 'Errors';
}

@Injectable({
  providedIn: 'root'
})
export class ConfigLogService {

  loggingLevel: LoggingLevel = LoggingLevel.None;
  constructor() { }
  
}
