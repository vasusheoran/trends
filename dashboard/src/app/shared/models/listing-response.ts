
export interface IUpdateResponse {
    CP:            number;
    Date:          Date;
    HP:            number;
    LP:            number;
    bi:            number;
    bj:            number;
    bk:            number;
    OP:            number;
}

export class UpdateResponse {
    CP:            number;
    Date:          Date;
    HP:            number;
    LP:            number;
    bi:            number;
    bj:            number;
    bk:            number;
    OP:            number;

    constructor(data){
        this.CP = data['CP']
        this.Date = data['Date']
        this.HP = data['HP']
        this.LP = data['LP']
        this.bi = data['bi']
        this.bj = data['bj']
        this.bk = data['bk']
        this.OP = data['OP']
    }
}