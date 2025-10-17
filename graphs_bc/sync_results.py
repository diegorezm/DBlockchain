import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import os
import matplotlib.ticker as mticker

# ---- Settings ----
INPUT_FILE = "./sync_results.csv"
TEST_NAME = "SyncTest"
OUTPUT_DIR = "./out_graphs"
os.makedirs(OUTPUT_DIR, exist_ok=True)

# ---- Read CSV ----
df = pd.read_csv(INPUT_FILE)

# Convert sync_time_ms to float (just in case it's stored as str)
df["sync_time_ms"] = pd.to_numeric(df["sync_time_ms"], errors="coerce")

# Compute mean sync time
mean_time = df["sync_time_ms"].mean()
print(f"Average sync time: {mean_time:.2f} ms")

# ---- Plot ----
sns.set_theme(style="whitegrid")
plt.figure(figsize=(7, 4))

sns.scatterplot(x="round", y="sync_time_ms", data=df, color="#4C6EF5", s=70)
plt.axhline(
    mean_time, color="#FF6B6B", linestyle="--", label=f"Mean {mean_time:.2f} ms"
)

plt.title(f"{TEST_NAME} – Synchronization Time per Round")
plt.xlabel("Round")
plt.ylabel("Sync Time (ms)")
plt.gca().yaxis.set_major_formatter(mticker.StrMethodFormatter("{x:.2f}"))
plt.legend(loc="upper right")
plt.tight_layout()

# Save output
output_path = os.path.join(OUTPUT_DIR, f"{TEST_NAME}_scatter_mean_time_.jpg")
plt.savefig(output_path, dpi=150, format="jpg")
plt.close()

print(f"[+] Saved graph: {output_path}")
