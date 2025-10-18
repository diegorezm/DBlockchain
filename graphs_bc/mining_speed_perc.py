import pandas as pd

df = pd.read_csv("./mining_speed_metrics.csv")

avg = df.groupby("difficulty")["duration_ms"].mean().reset_index()
avg["pct_increase_from_prev"] = avg["duration_ms"].pct_change() * 100

inc_3to4 = avg.loc[avg["difficulty"] == 4, "pct_increase_from_prev"].values[0]
inc_4to5 = avg.loc[avg["difficulty"] == 5, "pct_increase_from_prev"].values[0]
inc_5to6 = avg.loc[avg["difficulty"] == 6, "pct_increase_from_prev"].values[0]

val3 = avg.loc[avg["difficulty"] == 3, "duration_ms"].values[0]
val4 = avg.loc[avg["difficulty"] == 4, "duration_ms"].values[0]
val5 = avg.loc[avg["difficulty"] == 5, "duration_ms"].values[0]
val6 = avg.loc[avg["difficulty"] == 6, "duration_ms"].values[0]

print(f"📈 Aumento de {inc_3to4:.2f}% entre as dificuldades 3 e 4.")
print(f"📈 Aumento de {inc_4to5:.2f}% entre as dificuldades 4 e 5.")
print(f"📈 Aumento de {inc_5to6:.2f}% entre as dificuldades 5 e 6.")
