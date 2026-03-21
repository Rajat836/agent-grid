"use client";

import {
  BookOpen,
  Hexagon,
  Network,
  Sparkles,
} from "lucide-react";
import Link from "next/link";
import { Config } from "@/base/config";

const navItems = [
  {
    label: "Ontology Agent",
    href: "/system-analysis",
    icon: Sparkles,
    active: false,
  },
  {
    label: "Documentation Agent",
    href: "/knowledge-graph",
    icon: Network,
    active: false,
  },
  {
    label: "Documentation",
    href: "/documentation",
    icon: BookOpen,
    active: true,
  },
];

function cn(...classes) {
  return classes.filter(Boolean).join(" ");
}

function HexLogo() {
  return (
    <div className="relative flex h-10 w-10 items-center justify-center">
      <div className="absolute inset-0 rounded-[14px] bg-[linear-gradient(135deg,rgba(167,139,250,0.35),rgba(34,211,238,0.22))] blur-md" />
      <div className="relative flex h-10 w-10 items-center justify-center rounded-[14px] border border-white/10 bg-white/5 backdrop-blur-xl">
        <Hexagon className="h-5 w-5 text-[#c4b5fd]" />
      </div>
    </div>
  );
}

export default function DocumentationLayout({ children }) {
  return (
    <>
      <style jsx global>{`
        @import url("https://fonts.googleapis.com/css2?family=DM+Sans:wght@400;500&family=JetBrains+Mono:wght@400&family=Syne:wght@700&display=swap");
        
        .glass-panel {
          background: rgba(255, 255, 255, 0.04);
          border: 1px solid rgba(255, 255, 255, 0.08);
          backdrop-filter: blur(12px);
        }
      `}</style>

      <div
        className="min-h-screen overflow-hidden bg-[#0a0a18] text-[#e2e8f0]"
        style={{ fontFamily: "'DM Sans', sans-serif" }}
      >
        <div className="pointer-events-none absolute inset-0 z-0 overflow-hidden">
          <div className="absolute -left-20 -top-24 h-[400px] w-[400px] rounded-full bg-[rgba(139,92,246,0.12)] blur-[80px]" />
          <div className="absolute bottom-[-90px] right-[-60px] h-[350px] w-[350px] rounded-full bg-[rgba(6,182,212,0.08)] blur-[80px]" />
          <div className="absolute right-[18%] top-[28%] h-[280px] w-[280px] rounded-full bg-[rgba(236,72,153,0.06)] blur-[80px]" />
        </div>

        <div className="relative z-[1] flex min-h-screen">
          <aside className="hidden w-[200px] shrink-0 border-r border-[rgba(255,255,255,0.06)] bg-[rgba(255,255,255,0.02)] px-4 py-5 md:flex md:flex-col">
            <div className="flex items-center gap-3">
              <HexLogo />
              <div>
                <div
                  className="text-lg text-[#f5f3ff]"
                  style={{ fontFamily: "'Syne', sans-serif", fontWeight: 700 }}
                >
                  Ontology
                </div>
              </div>
            </div>

            <nav className="mt-8 space-y-2">
              {navItems.map((item) => {
                const Icon = item.icon;
                return (
                  <Link
                    key={item.label}
                    href={item.href}
                    className={cn(
                      "flex items-center gap-3 rounded-[12px] border border-transparent px-3 py-3 text-sm text-[rgba(255,255,255,0.6)] transition",
                      item.active
                        ? "border-[rgba(167,139,250,0.12)] bg-[rgba(167,139,250,0.1)] text-[#f5f3ff]"
                        : "hover:border-white/6 hover:bg-white/4 hover:text-[#e2e8f0]",
                    )}
                    style={
                      item.active
                        ? { boxShadow: "inset 2px 0 0 #a78bfa" }
                        : undefined
                    }
                  >
                    <Icon className="h-4 w-4" />
                    <span>{item.label}</span>
                  </Link>
                );
              })}
            </nav>
            
            <div className="mt-auto flex items-center justify-between rounded-[12px] border border-white/8 bg-white/4 px-3 py-2">
              <div className="text-xs text-[rgba(255,255,255,0.35)]">
                v0.1.0
              </div>
              <div className="flex items-center gap-2 text-xs text-[#6ee7b7]">
                <span className="h-2 w-2 rounded-full bg-[#6ee7b7] animate-pulse" />
                online
              </div>
            </div>
          </aside>

          <main className="flex min-h-screen flex-1 flex-col px-4 py-4 md:px-6 md:py-5">
            <div className="glass-panel rounded-[14px] px-5 py-4">
              <div className="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
                <div>
                  <div
                    className="text-[28px] leading-none text-[#f5f3ff]"
                    style={{
                      fontFamily: "'Syne', sans-serif",
                      fontWeight: 700,
                    }}
                  >
                    Documentation
                  </div>
                  <div className="mt-2 text-sm text-[rgba(255,255,255,0.35)]">
                    Understand how the agent works, what it can query, and how to get the most out of it.
                  </div>
                </div>
              </div>
            </div>
            {children}
          </main>
        </div>
      </div>
    </>
  );
}
