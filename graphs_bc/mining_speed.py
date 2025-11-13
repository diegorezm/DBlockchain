import os
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import matplotlib.ticker as mticker

INPUT_FILE = "./hashrate_metrics.csv"
TEST_NAME = "HashRate"
OUTPUT_DIR = "./out_graphs"
os.makedirs(OUTPUT_DIR, exist_ok=True)

df = pd.read_csv(INPUT_FILE)

# Média por dificuldade
avg = df.groupby("difficulty")["hash_rate_hs"].mean().reset_index()

# Variação percentual em relação à dificuldade 1
avg["pct_change"] = (avg["hash_rate_hs"] / avg["hash_rate_hs"].iloc[0] - 1) * 100

# ======== Gráfico 1: aumento/diminuição percentual do hash rate ========
sns.set_theme(style="whitegrid")
plt.figure(figsize=(8, 5))
sns.barplot(x="difficulty", y="pct_change", data=avg, color="#15AABF")
plt.title(f"{TEST_NAME} – Percentage Change in Hash Rate")
plt.ylabel("Change (%)")
plt.xlabel("Difficulty")
plt.tight_layout()

out_perc = os.path.join(OUTPUT_DIR, f"{TEST_NAME}_bar_perc_.jpg")
plt.savefig(out_perc, dpi=150, format="jpg")
plt.close()
print(f"[+] Saved {out_perc}")

# ======== Gráfico 2: distribuição de hash rate por dificuldade ========
plt.figure(figsize=(8, 5))
sns.boxplot(x="difficulty", y="hash_rate_hs", data=df, color="#FF922B")
plt.title(f"{TEST_NAME} – Hash Rate by Difficulty")
plt.xlabel("Difficulty")
plt.ylabel("Hash Rate (H/s)")

# formatar ticks com separador de milhar
plt.gca().yaxis.set_major_formatter(mticker.StrMethodFormatter("{x:,.0f} H/s"))
plt.tight_layout()

out_rate = os.path.join(OUTPUT_DIR, f"{TEST_NAME}_boxplot_rate_.jpg")
plt.savefig(out_rate, dpi=150, format="jpg")
plt.close()
print(f"[+] Saved {out_rate}")

print("Hash rate graphs generated")
