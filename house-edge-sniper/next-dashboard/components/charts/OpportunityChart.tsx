"use client";

import { useMemo } from "react";

interface Opportunity {
  timestamp_ms: number;
}

interface OpportunityChartProps {
  opportunities: Opportunity[];
}

/**
 * 簡易「機會數量」趨勢圖：依時間分桶顯示數量（純 CSS + 長條）
 * 若需進階圖表可改為 recharts 或 visx
 */
export function OpportunityChart({ opportunities }: OpportunityChartProps) {
  const buckets = useMemo(() => {
    const now = Date.now();
    const windowMs = 60 * 60 * 1000; // 1 小時
    const bucketCount = 24;
    const step = windowMs / bucketCount;
    const counts = new Array(bucketCount).fill(0);
    for (const opp of opportunities) {
      const age = now - opp.timestamp_ms;
      if (age < 0 || age > windowMs) continue;
      const idx = Math.min(bucketCount - 1, Math.floor((windowMs - age) / step));
      counts[idx]++;
    }
    const max = Math.max(1, ...counts);
    return counts.map((c) => (c / max) * 100);
  }, [opportunities]);

  return (
    <div
      style={{
        display: "flex",
        alignItems: "flex-end",
        gap: "4px",
        height: "100%",
        padding: "8px 0",
      }}
      aria-label="機會數量趨勢"
    >
      {buckets.map((pct, i) => (
        <div
          key={i}
          style={{
            flex: 1,
            minWidth: 4,
            height: `${Math.max(4, pct)}%`,
            background: "var(--accent)",
            borderRadius: "2px 2px 0 0",
            opacity: 0.6 + (pct / 100) * 0.4,
          }}
          title={`${buckets.length - i} 區間`}
        />
      ))}
    </div>
  );
}
