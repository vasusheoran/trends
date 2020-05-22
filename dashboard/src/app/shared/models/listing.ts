export interface IListing {
    CompanyName: string;
    Series: string;
    SASSymbol: string;
    YahooSymbol: string;
}

export class Listing implements IListing { 
  
    public CompanyName: string;
    public Series:string;
    public SASSymbol:string;
    public YahooSymbol:string;

    constructor(data: Partial<Listing>){
        Object.assign(this, data);
    }
  }