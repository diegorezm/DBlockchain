import requests
import time
import csv
from datetime import datetime

NODE_A = "http://localhost:3000"
NODE_B = "http://localhost:3001"
ROUNDS = 40
CSV_FILE = "sync_results.csv"


def mine_block(node_url):
    try:
        res = requests.post(f"{node_url}/api/chain/mine", timeout=10)
        res.raise_for_status()
        return True
    except requests.RequestException as e:
        print(f"[x] Mining failed on {node_url}: {e}")
        return False


def replace_chain(node_url):
    start = time.perf_counter()
    try:
        res = requests.post(f"{node_url}/api/chain/replace", timeout=10)
        res.raise_for_status()
        elapsed = (time.perf_counter() - start) * 1000  # ms
        return elapsed
    except requests.RequestException as e:
        print(f"[x] ReplaceChain failed on {node_url}: {e}")
        return None


def main():
    print(f"🏁 Benchmarking {ROUNDS} synchronization rounds...")
    results = []

    for i in range(ROUNDS):
        print(f"\n--- Round {i + 1} ---")

        print("⛏️  Minerando o Node A...")
        if not mine_block(NODE_A):
            continue

        print("🔄 Sync Node B com o Node A...")
        elapsed = replace_chain(NODE_B)
        if elapsed is not None:
            print(f"✅ Sync completed in {elapsed:.2f} ms")
            results.append(
                {
                    "timestamp": datetime.now().isoformat(),
                    "round": i + 1,
                    "sync_time_ms": f"{elapsed:.2f}",
                }
            )
        else:
            results.append(
                {
                    "timestamp": datetime.now().isoformat(),
                    "round": i + 1,
                    "sync_time_ms": "error",
                }
            )

    with open(CSV_FILE, "w", newline="") as csvfile:
        writer = csv.DictWriter(
            csvfile, fieldnames=["timestamp", "round", "sync_time_ms"]
        )
        writer.writeheader()
        writer.writerows(results)

    print("\n📊 Benchmark complete!")
    print(f"Results saved to: {CSV_FILE}")


if __name__ == "__main__":
    main()
