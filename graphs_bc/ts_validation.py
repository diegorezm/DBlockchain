import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import os
import matplotlib.ticker as mticker

# ---- Settings ----
INPUT_FILE = "./ts_validation_metrics.csv"
TEST_NAME = "TransactionValidation"
OUTPUT_DIR = "./out_graphs"
os.makedirs(OUTPUT_DIR, exist_ok=True)

# ---- Read the CSV ----
df = pd.read_csv(INPUT_FILE)
df_rounded = df.round({"duration_ms": 3}).drop_duplicates(subset=["duration_ms"])

# Compute mean validation time
mean_time = df_rounded["duration_ms"].mean()
print(f"Average validation time: {mean_time:.4f} ms")

# ---- Plot ----
sns.set_theme(style="whitegrid")
plt.figure(figsize=(7, 3))

# Scatter of unique (rounded) values
sns.scatterplot(x="iteration", y="duration_ms", data=df_rounded, color="#4C6EF5", s=50)

# Mean line + shaded band
plt.axhline(
    mean_time, color="#FF6B6B", linestyle="--", label=f"Mean {mean_time:.3f} ms"
)
plt.fill_between(
    [0, df_rounded["iteration"].max()],
    mean_time - 0.01,
    mean_time + 0.01,
    color="#FF6B6B",
    alpha=0.1,
)

plt.title("Average Transaction Validation Time")
plt.xlabel("Iteration (summary view)")
plt.ylabel("Duration (ms)")
plt.gca().yaxis.set_major_formatter(mticker.StrMethodFormatter("{x:.3f}"))
plt.legend(loc="upper right")
plt.tight_layout()

# Save output
output_path = os.path.join(OUTPUT_DIR, f"{TEST_NAME}_scatter_mean_time_.jpg")
plt.savefig(output_path, dpi=150, format="jpg")
plt.close()

print(f"[+] Saved graph: {output_path}")
