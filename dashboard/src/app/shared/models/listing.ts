export interface IListing {
    Company: string;
    Series: string;
    SAS: string;
    Symbol: string;
}

export class Listing implements IListing { 
  
    public Company: string;
    public Series:string;
    public SAS:string;
    public Symbol:string;

    constructor(data: Partial<Listing>){
        Object.assign(this, data);
    }
  }