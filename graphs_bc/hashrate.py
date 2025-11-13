import os
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import matplotlib.ticker as mticker

# Caminho do arquivo CSV com os dados
ARQUIVO_ENTRADA = "./hashrate_metrics.csv"
PASTA_SAIDA = "./out_graficos"
NOME_TEXTO = "Taxa de Hash"

os.makedirs(PASTA_SAIDA, exist_ok=True)

# Lê o CSV
df = pd.read_csv(ARQUIVO_ENTRADA)

# Calcula média e desvio padrão do hash rate por dificuldade
media = df.groupby("difficulty")["hash_rate_hs"].agg(["mean", "std"]).reset_index()

sns.set_theme(style="whitegrid")

# ======== Gráfico de linha (média + desvio padrão) ========
plt.figure(figsize=(8, 5))
sns.lineplot(
    data=media,
    x="difficulty",
    y="mean",
    marker="o",
    color="#15AABF",
    linewidth=2,
)

plt.title(f"{NOME_TEXTO} Média por Dificuldade")
plt.xlabel("Dificuldade")
plt.ylabel("Taxa de Hash (H/s)")

# Força o eixo X a mostrar apenas os valores inteiros de dificuldade
plt.xticks(media["difficulty"].unique())

# Formata o eixo Y com separadores de milhar
plt.gca().yaxis.set_major_formatter(mticker.StrMethodFormatter("{x:,.0f} H/s"))

plt.tight_layout()

saida = os.path.join(PASTA_SAIDA, f"{NOME_TEXTO.lower().replace(' ', '_')}_linha.jpg")
plt.savefig(saida, dpi=150, format="jpg")
plt.close()

print(f"[+] Gráfico salvo em: {saida}")
print("Gráfico de linha gerado com sucesso!")
