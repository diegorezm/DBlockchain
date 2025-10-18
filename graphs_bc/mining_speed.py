import os
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import matplotlib.ticker as mticker

INPUT_FILE = "./mining_speed_metrics.csv"
TEST_NAME = "MiningSpeed"
OUTPUT_DIR = "./out_graphs"
os.makedirs(OUTPUT_DIR, exist_ok=True)

df = pd.read_csv(INPUT_FILE)

avg = df.groupby("difficulty")["duration_ms"].mean().reset_index()
avg["pct_increase"] = (avg["duration_ms"] / avg["duration_ms"].iloc[0] - 1) * 100

sns.set_theme(style="whitegrid")
plt.figure(figsize=(8, 5))
sns.barplot(x="difficulty", y="pct_increase", data=avg, color="#4C6EF5")
plt.title(f"{TEST_NAME} – Percentage Increase in Mining Time")
plt.yscale("log")
plt.ylabel("Increase (log scale, %)")
plt.xlabel("Difficulty")
# plt.ylabel("Increase (%)")
plt.tight_layout()

out_perc = os.path.join(OUTPUT_DIR, f"{TEST_NAME}_bar_perc_.jpg")
plt.savefig(out_perc, dpi=150, format="jpg")
plt.close()
print(f"[+] Saved {out_perc}")

plt.figure(figsize=(8, 5))
sns.boxplot(x="difficulty", y="duration_ms", data=df, color="#82C91E")
plt.title(f"{TEST_NAME} – Mining Time by Difficulty")
plt.xlabel("Difficulty")
plt.ylabel("Duration (ms)")
plt.gca().yaxis.set_major_formatter(mticker.StrMethodFormatter("{x:,.0f} ms"))
plt.tight_layout()

out_time = os.path.join(OUTPUT_DIR, f"{TEST_NAME}_boxplot_time_.jpg")
plt.savefig(out_time, dpi=150, format="jpg")
plt.close()

print("Graphs generated")
