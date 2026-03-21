"use client";

import { useEffect, useId, useRef, useState } from "react";
import Script from "next/script";

declare global {
  interface Window {
    mermaid?: {
      initialize: (config: Record<string, unknown>) => void;
      render: (
        id: string,
        chart: string,
      ) => Promise<{ svg: string; bindFunctions?: Element }>;
      contentLoaded?: () => void;
    };
  }
}

const mermaidConfig = {
  startOnLoad: false,
  securityLevel: "loose",
  theme: "base",
  themeVariables: {
    primaryColor: "#1f2937",
    primaryTextColor: "#e2e8f0",
    primaryBorderColor: "#8b5cf6",
    lineColor: "#22d3ee",
    secondaryColor: "#0f172a",
    tertiaryColor: "#111827",
    clusterBkg: "#0f172a",
    clusterBorder: "#334155",
    mainBkg: "#111827",
    nodeBorder: "#8b5cf6",
    textColor: "#e2e8f0",
    fontFamily: "DM Sans, sans-serif",
    background: "transparent",
  },
  flowchart: {
    useMaxWidth: true,
    htmlLabels: true,
    curve: "basis",
  },
};

let mermaidInitialized = false;
let mermaidLoader: Promise<void> | null = null;

function waitForMermaid() {
  if (typeof window === "undefined") {
    return Promise.resolve();
  }

  if (window.mermaid) {
    return Promise.resolve();
  }

  if (!mermaidLoader) {
    mermaidLoader = new Promise((resolve, reject) => {
      const existingScript = document.querySelector(
        'script[data-mermaid-script="true"]',
      ) as HTMLScriptElement | null;

      if (existingScript?.getAttribute("data-loaded") === "true") {
        resolve();
        return;
      }

      const script =
        existingScript || document.createElement("script");

      script.src = "/vendor/mermaid.min.js";
      script.async = true;
      script.setAttribute("data-mermaid-script", "true");

      script.onload = () => {
        script.setAttribute("data-loaded", "true");
        resolve();
      };
      script.onerror = () => reject(new Error("Failed to load Mermaid."));

      if (!existingScript) {
        document.head.appendChild(script);
      }
    });
  }

  return mermaidLoader;
}

function initializeMermaid() {
  if (!window.mermaid || mermaidInitialized) {
    return;
  }

  window.mermaid.initialize(mermaidConfig);
  mermaidInitialized = true;
}

type MermaidDiagramProps = {
  chart: string;
  className?: string;
};

export function MermaidDiagram({
  chart,
  className = "",
}: MermaidDiagramProps) {
  const generatedId = useId().replace(/:/g, "");
  const containerRef = useRef<HTMLDivElement>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;

    async function renderChart() {
      if (!containerRef.current) {
        return;
      }

      await waitForMermaid();
      initializeMermaid();
      setError(null);

      try {
        const renderId = `mermaid-${generatedId}-${Date.now()}`;
        const { svg } = await window.mermaid.render(renderId, chart.trim());

        if (!active || !containerRef.current) {
          return;
        }

        containerRef.current.innerHTML = svg;
      } catch (renderError) {
        if (!active || !containerRef.current) {
          return;
        }

        containerRef.current.innerHTML = "";
        setError(
          renderError instanceof Error
            ? renderError.message
            : "Failed to render Mermaid chart.",
        );
      }
    }

    renderChart();

    return () => {
      active = false;
    };
  }, [chart, generatedId]);

  if (error) {
    return (
      <div
        className={`mermaid-error my-4 rounded-xl border border-rose-400/30 bg-rose-950/30 p-4 ${className}`.trim()}
      >
        <div className="mb-2 text-xs font-semibold uppercase tracking-[0.2em] text-rose-300">
          Mermaid render error
        </div>
        <pre className="overflow-x-auto whitespace-pre-wrap text-xs leading-6 text-rose-100">
          <code>{chart}</code>
        </pre>
        <div className="mt-3 text-xs text-rose-200/80">{error}</div>
      </div>
    );
  }

  return (
    <>
      <Script src="/vendor/mermaid.min.js" strategy="afterInteractive" />
      <div
        className={`mermaid-shell my-4 overflow-x-auto rounded-xl border border-white/10 bg-[radial-gradient(circle_at_top,rgba(139,92,246,0.12),transparent_45%),rgba(2,6,23,0.72)] p-4 ${className}`.trim()}
      >
        <div
          ref={containerRef}
          className="mermaid-diagram flex min-w-fit justify-center [&_svg]:h-auto [&_svg]:max-w-none"
        />
      </div>
    </>
  );
}
