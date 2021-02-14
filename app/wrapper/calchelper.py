# -*- coding: utf-8 -*-
"""
Created on Sun Mar 29 17:10:15 2020

@author: vsheoran
"""
from .utilities import Utilities
values = dict()

def ema_future(n, span ,df):
    last_row = df.iloc[0].to_dict()
    prev_ema = last_row['emaCP' + str(span)]
    
#        print(str(prev_ema) +  ' :: ' + str(last_row['CP']))
    for i in range(n):   
        prev_ema = ((2/(span+1)) * (last_row['CP'] - prev_ema)) + prev_ema
#        print(str(prev_ema))
    return prev_ema

#def find_ema(spans, cols ,df):
#    util = Utilities()
#    for span in spans:
#        df = util.ema_rolling(df, cols, span)
#    return df

def find_BX(df):
#        temp1 = last_day['HP']*2
    util = Utilities()
    HP =df.at[2,'HP']
    DEF = util.ema_calculate(5, df.at[2,'HP'],
                                  df.at[0, 'emaCP5'])
    GHI = util.ema_calculate(20, df.at[2,'HP'],
                                  df.at[0, 'emaCP20'])
    
    final=(HP+HP+(((((DEF)+((DEF)+((GHI)-(DEF))/2))/2)+((HP+(HP+((((DEF)+((DEF)+((GHI)-(DEF))/2))/2)-HP)/2))/2))/2))/3
    
    global values
    values.update({'bx' : final})
    
    return final
    

def find_CJ(df):    
    future_ema5 = ema_future(2,5, df)
    future_ema20 = ema_future(2,20,df)
    cj=(df.at[0,'CP']+df.at[0,'CP']+((((future_ema5+(future_ema5+(future_ema20-future_ema5)/2))/2)+((df.at[0,'CP']+(df.at[0,'CP']+(((future_ema5+(future_ema5+(future_ema20-future_ema5)/2))/2)-df.at[0,'CP'])/2))/2))/2))/3
    
    global values
    values.update({'cj' : cj})
    return cj

def find_U(df):        
    cj = find_CJ(df)
    cq = df.at[0,'minHP3'] - cj
    u = df.at[0,'minHP3'] + cq
#    db.update({'u' : {'cq' : cq , 'cj' : cj, 'val' : u}})
    global values
    values.update({'u' : u})
    return u

def find_co(row ,df):
    # cp_av = df[['CP_CI_HP']][row:row+10].values.mean(axis=0) + df[['CP_CI_HP']][row:row+50].values.mean(axis=0)
    # cv_av_1 = cp_av[0]
    # co = (cv_av_1)/2-((cv_av_1)/2*(((cv_av_1)/2-(((((cv_av_1)/2-((cv_av_1)/2*0.01))+(((cv_av_1)/2-((cv_av_1)/2*0.01))*0.025))+(cv_av_1)/2)/2))/(cv_av_1)/2*100/2)/100)
    cv_av_1 = df.at[row,'avCP_CI_HP10'] + df.at[row,'avCP_CI_HP50']
    co = (cv_av_1)/2-((cv_av_1)/2*(((cv_av_1)/2-(((((cv_av_1)/2-((cv_av_1)/2*0.01))+(((cv_av_1)/2-((cv_av_1)/2*0.01))*0.025))+(cv_av_1)/2)/2))/(cv_av_1)/2*100/2)/100)
    return co

def find_cp(row ,df):
    cp = df.at[row,'emaCP_CI_HP5']
    global values
    values.update({'cp' + str(row): cp})
    return cp
    
def find_ae(row ,df):
    ae_dict = dict()
    ae_dict.update({'co' : {'1' : find_co(row, df) , '2': find_co(row+1, df)},
                 'cp' : {'1' : find_cp(row, df) , '2': find_cp(row+1, df)}
                 })
    ae=df.at[row+1,'HP']-((ae_dict['cp']['1']-ae_dict['co']['1'])-(ae_dict['cp']['2']-ae_dict['co']['2']))
    global values
    ae_dict.update({'value' : ae})
    values.update({'ae-'+str(row) : ae_dict})
    return ae


def find_ai_af(df):
    u = find_U(df)
    bx = find_BX(df)            
    q = df.at[2,'HP']*2 - bx     
    # min_HP = min(df.at[2,'HP'], df.at[3,'HP'])
    min_HP = min(df.at[3,'HP'], df.at[2,'HP'])    
    ai=(( df.at[2,'LP'] + ( df.at[2,'HP'] + (u - df.at[2,'HP'] )/2 + min_HP + ((q) - min_HP)/2)/2)/2)
        
    af = df.at[2,'LP'] + ( ai-df.at[2,'LP'])/2   
     
    global values    
    values.update({'ai' : ai,'af' : af,'q' : q})
    
    return ai,af, u

def find_ao(row ,df):
    return find_ae(row, df) - find_ae(row+1, df)  
    

def find_BI(db, frozen, df , freeze = False):  
    global values 
    ai, af, u = find_ai_af(df)   
    ao = find_ao(1, df)    
    temp1 = af-ai
    temp2 = (ai+(af+(temp1)/2))/2
    temp3 = ((df.at[2,'LP']-ao)+(af+(temp1)/2))/2
    
    bi = max(temp2,temp3)
    
    if freeze:
        return bi
    
    frozen_bi = bi
    
    if 'bi' in frozen and 'Date' in frozen:
        frozen_bi = frozen['bi']
    
    find_BK(frozen_bi, df, bi)
    
    # find_AR(df)
    find_CR(df)

    # Trend
    bn, ar2 = find_BN_Row(df, 2)
    values.update({'trend': df.at[2,'CP'] - bn, 'ar': ar2})

    # Moment
    bn, ar2 = find_BN_Row(df, 1)
    values.update({'moment': df.at[2,'CP'] - bn})

    # Moment2
    # bn, ar2 = find_BN_Row(df, 3)
    # values.update({'moment2': df.at[2,'CP'] - bn})
    
    minHP2 = min(df.at[3,'HP'], df.at[4,'HP'])
       
    values.update({'bi' : bi, 
                   'frozen.bi': frozen_bi, 
                   'ao': ao, 
                   'min.HP.2' : minHP2,
                   'min.HP.3' : min(minHP2, df.at[5,'HP']),
                   'min.HP.f' : min(df.at[3,'HP'], df.at[2,'HP']),
                   'cr': find_CR(df)})
    
    return values

def find_BK(old_bi ,df, bi):
    bk = old_bi+(df.at[2,'HP']-old_bi)/2+(df.at[2,'HP']-(old_bi+(df.at[2,'HP']-old_bi)/2))/2
    bj = (bi + bk)/2 
    global values   
    values.update({'bk' : bk, 'bj' : bj})
    
def find_AR(df, row):
    # global values     
    # row = 2
    cp_av_10 = df[['CP']][row:row+10].values.mean(axis=0)
    cp_av_50 = df[['CP']][row:row+50].values.mean(axis=0)
    av = cp_av_10 + cp_av_50    # Sum
    ar = ((av)/2)-((av)/2*(((av)/2-(((((av)/2-((av)/2*0.01))+(((av)/2-((av)/2*0.01))*0.025))+(av)/2)/2))/(av)/2*100/2)/100)
    # values.update({'ar' : float(ar)})
    return ar

def find_CR(df):    
    # global values 
    ema_d = df.at[2, 'ema_diffCP1Pos']
    ema_e = df.at[2, 'ema_diffCP1Neg']
    
    if ema_e == 0:
        return 100.0
        # values.update({'cr' : 100.0})
        # return
    else:
        val = 100 - (100/(1+(ema_d)/ema_e))
        return val
        # values.update({'cr' : val})
        # return
    
# def find_BN(df):    
#     global values    
#     ar2 = float(find_AR(df, 2))
#     ar3 = float(find_AR(df, 3))
    
#     bn = (df.at[2,'emaCP5'] - ar2) - (df.at[3,'emaCP5'] - ar3)
#     values.update({'bn' : bn, 'ar' : ar2, 'trend': df.at[2,'CP'] - bn, 'Close': df.at[2,'CP']})
#     return bn, ar2
    
def find_BN_Row(df, row):    
    global values   
    first_index = row
    second_index = row + 1
    ar2 = float(find_AR(df, first_index))
    ar3 = float(find_AR(df, second_index))
    
    bn = (df.at[first_index,'emaCP5'] - ar2) - (df.at[second_index,'emaCP5'] - ar3)
    return bn, ar2
    
    