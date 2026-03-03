import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "House Edge Sniper | 戰情面板",
  description: "On-chain Casino Arb Bot 監控與設定",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="zh-TW">
      <body>
        <nav className="nav">
          <a href="/dashboard">戰情面板</a>
          <a href="/settings">設定</a>
        </nav>
        <main>{children}</main>
      </body>
    </html>
  );
}
