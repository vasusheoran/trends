
from wrapper.calculate import Calculate
from wrapper.database import DB
import pandas as pd


data = pd.read_csv("data/historical/Nifty 50/02-14-2021.csv")


calc = Calculate()
calc.process_file_or_df(data, "Nifty 50")

df = calc.get_dataframe()
for col in df.columns:
    print(col, end =" ")
print(col, end ="\n")

# print(df[["ema_diffCP1Pos","ema_diffCP1Neg", "CP", "diffCP1", "diffCP1Pos", "diffCP1Neg"]].head(10))
print(df[["CP5", "CP20", "emaCP5", "diffCP1Pos", "diffCP1Neg", "emaCP20", "emaCP_CI_HP20", "ema_diffCP1Pos", "ema_diffCP1Neg"]].head(6))


# val = calc.update({
#     'CP': 15163.299805,
#     'HP': 15243.50,
#     'LP': 15081.000000,
#     'Date': "",
# })
# print(calc.back_ground)
# print()
# print(val)
# print(df[["CP","CP_CI_HP"]].head(20))

