"use client";

import { useEffect, useState } from "react";
import { OpportunityChart } from "@/components/charts/OpportunityChart";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/ws";

interface Opportunity {
  chain: string;
  tx_hash: string;
  contract_address?: string;
  expected_value_wei?: string;
  timestamp_ms: number;
}

export default function DashboardPage() {
  const [opportunities, setOpportunities] = useState<Opportunity[]>([]);
  const [wsConnected, setWsConnected] = useState(false);
  const [risk, setRisk] = useState<{ max_drawdown_pct?: number }>({});

  useEffect(() => {
    const ws = new WebSocket(WS_URL);
    ws.onopen = () => setWsConnected(true);
    ws.onclose = () => setWsConnected(false);
    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data as string);
        if (data.type === "opportunity" && data.payload) {
          setOpportunities((prev) => [data.payload, ...prev].slice(0, 200));
        }
      } catch {
        // ignore
      }
    };
    return () => ws.close();
  }, []);

  useEffect(() => {
    async function fetchRisk() {
      try {
        const res = await fetch(`${API_BASE}/api/risk`);
        if (res.ok) setRisk(await res.json());
      } catch {
        // ignore
      }
    }
    fetchRisk();
  }, []);

  return (
    <div>
      <h1>戰情面板</h1>

      <div className="card">
        <span className={wsConnected ? "live-dot" : ""} />
        WebSocket: {wsConnected ? "已連線" : "未連線"}
        {risk.max_drawdown_pct != null && (
          <span style={{ marginLeft: "1.5rem", color: "var(--muted)" }}>
            最大回撤: {risk.max_drawdown_pct}%
          </span>
        )}
      </div>

      <div className="card">
        <h2 style={{ fontSize: "1.1rem", marginBottom: "1rem" }}>機會數量趨勢</h2>
        <div className="chart-container">
          <OpportunityChart opportunities={opportunities} />
        </div>
      </div>

      <div className="card">
        <h2 style={{ fontSize: "1.1rem", marginBottom: "1rem" }}>即時機會</h2>
        <table>
          <thead>
            <tr>
              <th>鏈</th>
              <th>Tx Hash</th>
              <th>合約</th>
              <th>期望值 (wei)</th>
              <th>時間</th>
            </tr>
          </thead>
          <tbody>
            {opportunities.slice(0, 20).map((opp, i) => (
              <tr key={`${opp.tx_hash}-${i}`}>
                <td>{opp.chain}</td>
                <td>
                  <code style={{ fontSize: "0.85rem" }}>
                    {opp.tx_hash.slice(0, 18)}…
                  </code>
                </td>
                <td>{opp.contract_address || "—"}</td>
                <td>{opp.expected_value_wei ?? "—"}</td>
                <td>
                  {new Date(opp.timestamp_ms).toLocaleTimeString("zh-TW")}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {opportunities.length === 0 && (
          <p style={{ color: "var(--muted)", padding: "1rem 0" }}>
            尚無機會資料，請確認 Rust Sniper 與 Go Coordinator 已啟動。
          </p>
        )}
      </div>
    </div>
  );
}
