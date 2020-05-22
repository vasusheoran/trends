export interface IRealTimeDataResponse {
    CP:     number;
    Date:   Date;
}

export class RealTimeDataResponse {
    CP:     number;
    Date:   Date;

    constructor(data){
        this.CP = data['CP']
        this.Date = data['Date']
    }
}