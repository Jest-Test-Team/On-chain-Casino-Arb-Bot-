"use client";

import { useState, useEffect } from "react";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export default function SettingsPage() {
  const [wallet, setWallet] = useState("");
  const [minBalanceWei, setMinBalanceWei] = useState("");
  const [maxDrawdownPct, setMaxDrawdownPct] = useState("20");
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    // 可從 API 載入既有設定
    async function load() {
      try {
        const res = await fetch(`${API_BASE}/api/risk`);
        if (res.ok) {
          const data = await res.json();
          if (data.max_drawdown_pct != null)
            setMaxDrawdownPct(String(data.max_drawdown_pct));
        }
      } catch {
        // ignore
      }
    }
    load();
  }, []);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    try {
      const res = await fetch(`${API_BASE}/api/settings`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          wallet: wallet || undefined,
          min_balance_wei: minBalanceWei || undefined,
          max_drawdown_pct: maxDrawdownPct ? parseFloat(maxDrawdownPct) : undefined,
        }),
      });
      if (res.ok) setSaved(true);
    } catch {
      setSaved(false);
    }
  }

  return (
    <div>
      <h1>設定</h1>

      <div className="card">
        <h2 style={{ fontSize: "1.1rem", marginBottom: "1rem" }}>
          錢包與策略參數
        </h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>錢包位址（選填）</label>
            <input
              type="text"
              placeholder="0x..."
              value={wallet}
              onChange={(e) => setWallet(e.target.value)}
            />
          </div>
          <div className="form-group">
            <label>最低餘額 (wei)</label>
            <input
              type="text"
              placeholder="例如 1000000000000000000"
              value={minBalanceWei}
              onChange={(e) => setMinBalanceWei(e.target.value)}
            />
          </div>
          <div className="form-group">
            <label>最大回撤 (%)</label>
            <input
              type="number"
              min="0"
              max="100"
              step="0.5"
              value={maxDrawdownPct}
              onChange={(e) => setMaxDrawdownPct(e.target.value)}
            />
          </div>
          <button type="submit" className="btn">
            儲存
          </button>
          {saved && (
            <span style={{ marginLeft: "1rem", color: "var(--accent)" }}>
              已儲存
            </span>
          )}
        </form>
        <p style={{ color: "var(--muted)", fontSize: "0.9rem", marginTop: "1rem" }}>
          註：設定 API 需在 Go Coordinator 實作 POST /api/settings 後才會寫入。
        </p>
      </div>
    </div>
  );
}
