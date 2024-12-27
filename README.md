# trends
Go based app to visualize market trends.

Low:			If current day Low (23685) < previous day Low (23647) then Low colour dark:text-red-200 else ""
			
H/L:			If current day Close (23728) > current day H/l (23870) then H/L colour  dark:text-red-200 else ""
			If current day Close (23728) < current day H/l (23870) then H/L colour this
			
AVG & EMA-5:			If current day EMA-5 (23874) > current day AVG (24270) then AVG & EMA-5Background colour is  dark:text-purple-200 
			If current day EMA-5 (23874) < current day AVG (24270) then AVG & EMA-5Background colour is  dark:text-pink-200 
			
AVG:			If current day Close (23728) > current day AVG (24270) then AVG colour is   dark:text-purple-200 
			If current day Close (23728) < current day AVG (24270) then AVG colour is dark:text-pink-200 
			
EMA-5:			If current day Close (23728) > current day EMA-5 (23874) then EMA-5 colour is dark:text-pink-200 
			If current day Close (23728) < current day EMA-5 (23874) then EMA-5 colour is dark:text-red-200
			
EMA-20:			If current day Close (23728) > current day EMA-20 (24186) then EMA-20 colour is dark:text-blue-200 
			If current day Close (23728) < current day EMA-20 (24186) then EMA-20 colour is dark:text-red-200 
			
EMA:			If current day EMA (-312) > previous day EMA (-287) then EMA colour is dark:text-blue-200 
			If current day EMA (-312) < previous day EMA (-287) then EMA colour is dark:text-red-200 
			
Buy:			If current day Close (23728) > current day Buy (23911) then buy colour is dark:text-purple-200 
			If current day Close (23728) < current day Buy (23911) then buy colour is dark:text-red-200 
			If current day Close (23728) > current day Buy (23911) But current day low (23685) > current day Buy(23911) then colour is dark:text-red-200 
			
Support:			If current day Close (23728) > current day Support (24013) then buy colour is dark:text-purple-200 
			If current day Close (23728) < current day Support (24013) then Support colour is dark:text-red-200 
			If current day Close (23728) > current day Support (24013) But current day low (23685) < current day Support(24013) then colour is dark:text-green-200 
			If current day Close (23728) > current day Support (24013) But current day low (23685) > current day Support(24013) then colour is dark:text-red-200 
			
SMA:			If current day Close (23728) > current day SMA (23765) then buy colour is dark:text-purple-200 
			If current day Close (23728) < current day SMA (23765) then SMA colour is dark:text-red-200 
			
Buy, SMA and Suppport:	When current close (23728) > current day Buy,Support & SMA(all three) then background colour is dark:text-purple-200 

RSI:			If current day RSI (38) < 50.00 then RSI colour is dark:text-red-200 
			If current day RSI (38) > 50.00 But < 60.00 then RSI colour is black
			If current day RSI (38) > 60.00 But < 70.00 then RSI colour is dark:text-green-200 
			If current day RSI (38) > 70.00 then RSI colour is dark:text-blue-200 
			
CH:			If current day Close (23728) > current day CH (23912) then CH colour is dark:text-blue-200 
			If current day High (23868) < current day CH (23912) then CH colour is pink