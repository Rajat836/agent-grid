"use client";

import {
  BookOpen,
  Check,
  Copy,
  Hexagon,
  Lock,
  Mail,
  Network,
  Send,
  Sparkles,
} from "lucide-react";
import Link from "next/link";
import { useEffect, useMemo, useRef, useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { Config } from "@/base/config";

const samplePrompts = [
  "Give me list of services",
  "Give me list of teams",
  "List all entities under customer onboarding flow",
  "Give me list of api's under set password entity",
];

const navItems = [
  {
    label: "Ontology Agent",
    href: "/system-analysis",
    icon: Sparkles,
    active: true,
  },
  {
    label: "Knowledge Graph",
    href: "/knowledge-graph",
    icon: Network,
    active: false,
  },
  {
    label: "Documentation",
    href: "/documentation",
    icon: BookOpen,
    active: false,
  },
];

function cn(...classes) {
  return classes.filter(Boolean).join(" ");
}

function getSummaryWebSocketUrl() {
  const baseUrl =
    Config.apiHost ||
    (typeof window !== "undefined" ? window.location.origin : "");

  if (!baseUrl) {
    return "";
  }

  const url = new URL("/agent/v1/ontology/summary", baseUrl);
  url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
  return url.toString();
}

function formatEmailSection(title, content) {
  return `${title}\n${"=".repeat(title.length)}\n${content}`.trim();
}

function formatEmailPayload(payload) {
  if (!payload) {
    return "No response body was returned for this analysis.";
  }

  if (typeof payload !== "string") {
    return JSON.stringify(payload, null, 2);
  }

  const trimmedPayload = payload.trim();
  if (!trimmedPayload) {
    return "No response body was returned for this analysis.";
  }

  try {
    return JSON.stringify(JSON.parse(trimmedPayload), null, 2);
  } catch {
    return trimmedPayload;
  }
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

function MarkdownContent({ content }) {
  return (
    <div className="ontology-markdown max-w-none text-[#e2e8f0]">
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={{
          h2: ({ ...props }) => <h2 {...props} />,
          h3: ({ ...props }) => <h3 {...props} />,
          p: ({ ...props }) => <p {...props} />,
          ul: ({ ...props }) => <ul {...props} />,
          ol: ({ ...props }) => <ol {...props} />,
          li: ({ ...props }) => <li {...props} />,
          table: ({ ...props }) => (
            <div className="overflow-x-auto">
              <table {...props} />
            </div>
          ),
          code: ({ inline, children, ...props }) =>
            inline
              ? <code {...props}>{children}</code>
              : <pre>
                  <code {...props}>{children}</code>
                </pre>,
          strong: ({ ...props }) => <strong {...props} />,
          em: ({ ...props }) => <em {...props} />,
        }}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
}

function ShimmerCard() {
  return (
    <div className="glass-panel relative overflow-hidden rounded-[14px]">
      <div className="h-[2px] w-full bg-[linear-gradient(90deg,transparent,rgba(167,139,250,0.95),rgba(34,211,238,0.95),transparent)]" />
      <div className="p-5">
        <div className="flex items-center justify-between gap-3">
          <div className="flex items-center gap-3">
            <div className="h-8 w-8 rounded-full border border-white/10 bg-white/5">
              <div className="mx-auto mt-2 h-4 w-4 rounded-full border-2 border-transparent border-t-[#a78bfa] animate-spin" />
            </div>
            <div>
              <div className="text-[10px] uppercase tracking-[0.24em] text-[rgba(255,255,255,0.35)]">
                Ontology Agent
              </div>
              <div className="text-sm font-medium text-[#e2e8f0]">
                Generating report
              </div>
            </div>
          </div>
          <div className="rounded-full border border-white/10 bg-white/5 px-2.5 py-1 text-[10px] text-[rgba(255,255,255,0.35)]">
            last_update pending
          </div>
        </div>

        <div className="mt-6 space-y-3">
          {[90, 76, 82, 68, 54].map((width, index) => (
            <div
              key={width}
              className="h-3 rounded-full bg-white/10 shimmer-bar"
              style={{ width: `${width}%`, animationDelay: `${index * 0.15}s` }}
            />
          ))}
        </div>
      </div>
    </div>
  );
}

function RingsEmptyState() {
  return (
    <div className="flex min-h-[260px] flex-col items-center justify-center px-6 text-center">
      <div className="relative flex h-28 w-28 items-center justify-center">
        <span className="ring ring-1" />
        <span className="ring ring-2" />
        <span className="ring ring-3" />
        <span className="relative z-10 h-3 w-3 rounded-full bg-[#a78bfa] shadow-[0_0_22px_rgba(167,139,250,0.7)]" />
      </div>
      <div className="mt-6 text-sm text-[rgba(255,255,255,0.35)]">
        Ask Ontology Agent about services, teams, features, entities, and APIs.
      </div>
    </div>
  );
}

function formatPayload(payload) {
  return String(payload || "")
    .split("\n")
    .map((line) => `$ ${line || " "}`)
    .join("\n");
}

function StepCard({ step, index, state, onToggle }) {
  return (
    <div
      className="relative pl-10 opacity-0 animate-[fadeInUp_0.4s_ease-out_forwards]"
      style={{ animationDelay: `${index * 120}ms` }}
    >
      <div className="absolute left-[15px] top-0 h-full w-[2px] rounded-full bg-[linear-gradient(180deg,rgba(167,139,250,0.7),rgba(34,211,238,0.35))]" />
      <div className="absolute left-0 top-5 z-10 flex h-8 w-8 items-center justify-center rounded-full border border-white/8 bg-[#0f1020]">
        {state === "done"
          ? <span className="flex h-6 w-6 items-center justify-center rounded-full bg-[rgba(110,231,183,0.14)] text-[#6ee7b7]">
              <Check className="h-3.5 w-3.5" />
            </span>
          : state === "active"
            ? <span className="h-4 w-4 rounded-full border-2 border-transparent border-t-[#a78bfa] animate-spin" />
            : <span className="h-2.5 w-2.5 rounded-full bg-[rgba(255,255,255,0.24)]" />}
      </div>

      <button
        type="button"
        onClick={onToggle}
        className={cn(
          "glass-panel w-full rounded-[14px] px-4 py-4 text-left transition",
          state === "active" &&
            "border-[rgba(167,139,250,0.3)] shadow-[0_0_20px_rgba(139,92,246,0.08)]",
        )}
      >
        <div className="flex items-start justify-between gap-3">
          <div>
            <div
              className={cn(
                "text-[13px] font-medium",
                state === "active"
                  ? "text-[#e2e8f0]"
                  : "text-[rgba(255,255,255,0.4)]",
              )}
            >
              {step.message}
            </div>
            <div className="mt-1 text-[10px] uppercase tracking-[0.22em] text-[rgba(255,255,255,0.2)]">
              {state === "active" ? "In progress" : "Captured"}
            </div>
          </div>
          <div className="font-['JetBrains_Mono',monospace] text-[10px] text-[rgba(255,255,255,0.2)]">
            {step.elapsedMs != null
              ? `${(step.elapsedMs / 1000).toFixed(2)}s`
              : "--"}
          </div>
        </div>

        {step.expanded
          ? <pre className="mt-4 overflow-x-auto rounded-[12px] border border-white/6 bg-[rgba(0,0,0,0.4)] p-3 font-['JetBrains_Mono',monospace] text-[11px] leading-5 text-[rgba(255,255,255,0.45)]">
              <code>{formatPayload(step.payload)}</code>
            </pre>
          : null}
      </button>
    </div>
  );
}

export default function SystemAnalysisPage() {
  const [prompt, setPrompt] = useState("");
  const [runs, setRuns] = useState([]);
  const [steps, setSteps] = useState([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [activeRunId, setActiveRunId] = useState(null);
  const [analysisStartedAt, setAnalysisStartedAt] = useState(null);
  const [analysisFinishedAt, setAnalysisFinishedAt] = useState(null);
  const [copiedRunId, setCopiedRunId] = useState(null);
  const [emailFeedback, setEmailFeedback] = useState({});
  const [now, setNow] = useState(Date.now());
  const socketRef = useRef(null);
  const chatScrollRef = useRef(null);

  useEffect(() => {
    const interval = window.setInterval(() => {
      setNow(Date.now());
    }, 100);

    return () => window.clearInterval(interval);
  }, []);

  useEffect(() => {
    if (chatScrollRef.current) {
      chatScrollRef.current.scrollTop = chatScrollRef.current.scrollHeight;
    }
  }, []);

  useEffect(() => {
    return () => {
      if (socketRef.current) {
        socketRef.current.close();
      }
    };
  }, []);

  const totalElapsedMs = useMemo(() => {
    if (!analysisStartedAt) {
      return 0;
    }

    return Math.max(0, (analysisFinishedAt ?? now) - analysisStartedAt);
  }, [analysisFinishedAt, analysisStartedAt, now]);

  const finalizeRun = (finishedAt) => {
    setSteps((currentSteps) =>
      currentSteps.map((step) => ({
        ...step,
        status: "done",
        expanded: false,
        endedAt: finishedAt,
        elapsedMs: step.elapsedMs ?? finishedAt - step.startedAt,
      })),
    );
    setIsStreaming(false);
    setAnalysisFinishedAt(finishedAt);
  };

  const setRunState = (runId, updater) => {
    setRuns((currentRuns) =>
      currentRuns.map((run) =>
        run.id === runId ? { ...run, ...updater(run) } : run,
      ),
    );
  };

  const getEmailBody = (run) => {
    const responseBody =
      run.response?.payload ||
      run.error?.payload ||
      "No response body was returned for this analysis.";

    const sections = [
      formatEmailSection(
        "Ontology Agent Analysis",
        `Generated at: ${new Date(run.finishedAt || run.startedAt).toLocaleString()}`,
      ),
      formatEmailSection("Prompt", run.prompt.trim()),
      formatEmailSection("Response", formatEmailPayload(responseBody)),
    ];

    return sections.join("\n\n");
  };

  const startAnalysis = (submittedPrompt) => {
    if (socketRef.current) {
      socketRef.current.close();
      socketRef.current = null;
    }

    const websocketUrl = getSummaryWebSocketUrl();
    const startedAt = Date.now();
    const runId = `run-${startedAt}`;
    const run = {
      id: runId,
      prompt: submittedPrompt,
      startedAt,
      finishedAt: null,
      status: "streaming",
      response: null,
      error: null,
    };

    setRuns((currentRuns) => [...currentRuns, run]);
    setSteps([]);
    setCopiedRunId(null);
    setActiveRunId(runId);
    setIsStreaming(true);
    setAnalysisStartedAt(startedAt);
    setAnalysisFinishedAt(null);

    if (!websocketUrl) {
      const finishedAt = Date.now();
      setSteps([
        {
          id: `step-error-${finishedAt}`,
          type: "error",
          message: "Configuration error",
          payload:
            "Could not resolve the WebSocket host for /agent/v1/ontology/summary.",
          status: "done",
          startedAt,
          endedAt: finishedAt,
          elapsedMs: finishedAt - startedAt,
          expanded: true,
        },
      ]);
      setRunState(runId, () => ({
        finishedAt,
        status: "error",
        error: {
          message: "Configuration error",
          payload:
            "Could not resolve the WebSocket host for /agent/v1/ontology/summary.",
          createdAt: finishedAt,
        },
      }));
      setAnalysisFinishedAt(finishedAt);
      setIsStreaming(false);
      return;
    }

    const socket = new WebSocket(websocketUrl);
    socketRef.current = socket;

    socket.onopen = () => {
      socket.send(JSON.stringify({ prompt: submittedPrompt }));
    };

    socket.onmessage = (event) => {
      const eventTime = Date.now();

      try {
        const data = JSON.parse(event.data);

        if (data.type === "error") {
          setSteps((currentSteps) => {
            const finalized = currentSteps.map((step) =>
              step.status === "active"
                ? {
                    ...step,
                    status: "done",
                    expanded: false,
                    endedAt: eventTime,
                    elapsedMs: eventTime - step.startedAt,
                  }
                : step,
            );

            return [
              ...finalized,
              {
                id: `step-${eventTime}`,
                type: "error",
                message: data.message || "Processing error",
                payload:
                  typeof data.payload === "string"
                    ? data.payload
                    : data.payload != null
                      ? JSON.stringify(data.payload, null, 2)
                      : "",
                status: "done",
                startedAt: eventTime,
                endedAt: eventTime,
                elapsedMs: 0,
                expanded: true,
              },
            ];
          });

          setIsStreaming(false);
          setAnalysisFinishedAt(eventTime);
          setRunState(runId, () => ({
            finishedAt: eventTime,
            status: "error",
            error: {
              message: data.message || "Processing error",
              payload:
                typeof data.payload === "string"
                  ? data.payload
                  : data.payload != null
                    ? JSON.stringify(data.payload, null, 2)
                    : "",
              createdAt: eventTime,
            },
          }));
          socket.close();
          socketRef.current = null;
          return;
        }

        if (data.type === "last_update") {
          const result = {
            id: `last-${eventTime}`,
            type: "last_update",
            message: data.message || "Analysis Complete",
            payload:
              typeof data.payload === "string"
                ? data.payload
                : JSON.stringify(data.payload, null, 2),
            createdAt: eventTime,
          };

          finalizeRun(eventTime);
          setRunState(runId, () => ({
            finishedAt: eventTime,
            status: "done",
            response: result,
          }));
          socket.close();
          socketRef.current = null;
          return;
        }

        setSteps((currentSteps) => {
          const finalized = currentSteps.map((step) =>
            step.status === "active"
              ? {
                  ...step,
                  status: "done",
                  expanded: false,
                  endedAt: eventTime,
                  elapsedMs: eventTime - step.startedAt,
                }
              : step,
          );

          return [
            ...finalized,
            {
              id: `step-${eventTime}`,
              type: data.type || "update",
              message: data.message || "Processing",
              payload:
                typeof data.payload === "string"
                  ? data.payload
                  : data.payload != null
                    ? JSON.stringify(data.payload, null, 2)
                    : "",
              status: "active",
              startedAt: eventTime,
              elapsedMs: null,
              expanded: true,
            },
          ];
        });
      } catch {
        setSteps((currentSteps) => [
          ...currentSteps,
          {
            id: `step-parse-${eventTime}`,
            type: "error",
            message: "Invalid server message",
            payload: String(event.data || "Unknown event payload"),
            status: "done",
            startedAt: eventTime,
            endedAt: eventTime,
            elapsedMs: 0,
            expanded: true,
          },
        ]);
        setIsStreaming(false);
        setAnalysisFinishedAt(eventTime);
        setRunState(runId, () => ({
          finishedAt: eventTime,
          status: "error",
          error: {
            message: "Invalid server message",
            payload: String(event.data || "Unknown event payload"),
            createdAt: eventTime,
          },
        }));
        socket.close();
        socketRef.current = null;
      }
    };

    socket.onerror = (event) => {
      const errorTime = Date.now();

      setSteps((currentSteps) => [
        ...currentSteps.map((step) =>
          step.status === "active"
            ? {
                ...step,
                status: "done",
                expanded: false,
                endedAt: errorTime,
                elapsedMs: errorTime - step.startedAt,
              }
            : step,
        ),
        {
          id: `step-error-${errorTime}`,
          type: "error",
          message: "Connection error",
          payload: `Failed to connect to /agent/v1/ontology/summary.

Error details: ${event.message || event.type || "Unknown error"}
`,
          status: "done",
          startedAt: errorTime,
          endedAt: errorTime,
          elapsedMs: 0,
          expanded: true,
        },
      ]);

      setIsStreaming(false);
      setAnalysisFinishedAt(errorTime);
      setRunState(runId, () => ({
        finishedAt: errorTime,
        status: "error",
        error: {
          message: "Connection error",
          payload: `Failed to connect to /agent/v1/ontology/summary.

Error details: ${event.message || event.type || "Unknown error"}
`,
          createdAt: errorTime,
        },
      }));
    };

    socket.onclose = (event) => {
      if (socketRef.current === socket) {
        socketRef.current = null;
      }
      console.warn(
        "WebSocket connection closed:",
        `Code: ${event.code}`,
        `Reason: ${event.reason}`,
        `Was Clean: ${event.wasClean}`,
      );
    };
  };

  const handleSubmit = () => {
    const trimmed = prompt.trim();
    if (!trimmed || isStreaming) {
      return;
    }

    startAnalysis(trimmed);
  };

  const handlePromptKeyDown = (event) => {
    if (event.key === "Enter" && !event.shiftKey) {
      event.preventDefault();
      handleSubmit();
    }
  };

  const handleSuggestion = (value) => {
    setPrompt(value);
    if (!isStreaming) {
      startAnalysis(value);
    }
  };

  const handleCopy = async (run) => {
    if (!run?.response?.payload) {
      return;
    }

    try {
      await navigator.clipboard.writeText(run.response.payload);
      setCopiedRunId(run.id);
      window.setTimeout(
        () =>
          setCopiedRunId((current) => (current === run.id ? null : current)),
        1400,
      );
    } catch {
      setCopiedRunId(null);
    }
  };

  const handleEmailSend = (run) => {
    const subject = `Ontology Agent response · ${new Date(run.finishedAt || run.startedAt).toLocaleString()}`;
    const body = getEmailBody(run);
    const mailtoUrl = `mailto:?subject=${encodeURIComponent(subject)}&body=${encodeURIComponent(body)}`;

    window.location.assign(mailtoUrl);
    setEmailFeedback((current) => ({
      ...current,
      [run.id]: "Opened your email client.",
    }));
  };

  const activeRun = runs.find((run) => run.id === activeRunId) || null;

  const renderRunCard = (run) => {
    const resultEvent = run.response || run.error;
    const isCompleted = Boolean(run.response);
    const feedback = emailFeedback[run.id];

    return (
      <div key={run.id} className="mt-5 space-y-4">
        <div className="flex justify-end">
          <div className="max-w-2xl rounded-[14px_4px_14px_14px] bg-[linear-gradient(135deg,#5b21b6,#0e7490)] px-4 py-3 text-sm text-white shadow-[0_4px_24px_rgba(124,58,237,0.3)]">
            {run.prompt}
          </div>
        </div>

        {run.status === "streaming"
          ? <ShimmerCard />
          : resultEvent
            ? <div className="glass-panel overflow-hidden rounded-[14px] animate-[resultEnter_0.4s_ease-out_forwards]">
                <div
                  className={cn(
                    "h-[2px] w-full",
                    isCompleted
                      ? "bg-[linear-gradient(90deg,transparent,#a78bfa,#22d3ee,transparent)]"
                      : "bg-[linear-gradient(90deg,transparent,#fb7185,#f59e0b,transparent)]",
                  )}
                />
                <div className="p-5">
                  <div className="flex flex-col gap-3 border-b border-white/8 pb-4 lg:flex-row lg:items-center lg:justify-between">
                    <div className="flex items-center gap-3">
                      <div
                        className={cn(
                          "flex h-9 w-9 items-center justify-center rounded-full border border-white/8",
                          isCompleted
                            ? "bg-[rgba(167,139,250,0.08)] text-[#c4b5fd]"
                            : "bg-[rgba(251,113,133,0.08)] text-[#fda4af]",
                        )}
                      >
                        <Sparkles className="h-4 w-4" />
                      </div>
                      <div>
                        <div className="flex items-center gap-2">
                          <span
                            className="text-[15px] text-[#f5f3ff]"
                            style={{
                              fontFamily: "'Syne', sans-serif",
                              fontWeight: 700,
                            }}
                          >
                            {resultEvent.message ||
                              (isCompleted
                                ? "Analysis Complete"
                                : "Analysis Error")}
                          </span>
                          <span
                            className={cn(
                              "rounded-full px-2 py-0.5 text-[10px] uppercase tracking-[0.18em]",
                              isCompleted
                                ? "border border-[rgba(167,139,250,0.22)] bg-[rgba(167,139,250,0.08)] text-[#c4b5fd]"
                                : "border border-[rgba(251,113,133,0.22)] bg-[rgba(251,113,133,0.08)] text-[#fda4af]",
                            )}
                          >
                            {isCompleted ? "last_update" : "error"}
                          </span>
                        </div>
                        <div className="mt-1 text-xs text-[rgba(255,255,255,0.35)]">
                          {new Date(resultEvent.createdAt).toLocaleTimeString(
                            [],
                            {
                              hour: "2-digit",
                              minute: "2-digit",
                              second: "2-digit",
                            },
                          )}
                        </div>
                      </div>
                    </div>

                    {isCompleted
                      ? <div className="flex flex-wrap items-center justify-end gap-2 self-start">
                          <button
                            type="button"
                            onClick={() => handleEmailSend(run)}
                            className="inline-flex items-center gap-2 rounded-[10px] border border-[rgba(34,211,238,0.24)] bg-[rgba(34,211,238,0.08)] px-3 py-2 text-xs text-[#e2e8f0] transition hover:bg-[rgba(34,211,238,0.14)]"
                          >
                            <Mail className="h-3.5 w-3.5" />
                            Send Email
                          </button>
                          <button
                            type="button"
                            onClick={() => handleCopy(run)}
                            className="inline-flex items-center gap-2 rounded-[10px] border border-white/8 bg-white/4 px-3 py-2 text-xs text-[#e2e8f0] transition hover:border-[rgba(34,211,238,0.22)] hover:bg-[rgba(34,211,238,0.08)]"
                          >
                            <Copy className="h-3.5 w-3.5" />
                            {copiedRunId === run.id ? "Copied" : "Copy"}
                          </button>
                        </div>
                      : null}
                  </div>

                  <div className="mt-5">
                    <MarkdownContent content={resultEvent.payload} />
                  </div>

                  {isCompleted
                    ? <div className="mt-3">
                        {feedback
                          ? <div className="text-xs text-[rgba(255,255,255,0.55)]">
                              {feedback}
                            </div>
                          : null}
                      </div>
                    : null}
                </div>
              </div>
            : null}
      </div>
    );
  };

  const completionLabel =
    activeRun?.response && analysisFinishedAt
      ? `Analysis complete  ·  ${steps.length} steps  ·  ${(totalElapsedMs / 1000).toFixed(1)}s`
      : null;

  return (
    <>
      <style jsx global>{`
        @import url("https://fonts.googleapis.com/css2?family=DM+Sans:wght@400;500&family=JetBrains+Mono:wght@400&family=Syne:wght@700&display=swap");

        @keyframes shimmerPulse {
          0%,
          100% {
            opacity: 0.3;
          }
          50% {
            opacity: 0.6;
          }
        }

        @keyframes fadeInUp {
          from {
            opacity: 0;
            transform: translateY(8px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }

        @keyframes resultEnter {
          from {
            opacity: 0;
            transform: translateY(10px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }

        @keyframes ringPulse {
          0% {
            transform: scale(0.82);
            opacity: 0.08;
          }
          50% {
            opacity: 0.22;
          }
          100% {
            transform: scale(1.08);
            opacity: 0.04;
          }
        }

        .glass-panel {
          background: rgba(255, 255, 255, 0.04);
          border: 1px solid rgba(255, 255, 255, 0.08);
          backdrop-filter: blur(12px);
        }

        .shimmer-bar {
          animation: shimmerPulse 1.5s ease-in-out infinite;
        }

        .ring {
          position: absolute;
          inset: 0;
          border-radius: 9999px;
          border: 1px solid rgba(167, 139, 250, 0.2);
          animation: ringPulse 2.4s ease-in-out infinite;
        }

        .ring-2 {
          inset: 12px;
          animation-delay: 0.35s;
          border-color: rgba(34, 211, 238, 0.18);
        }

        .ring-3 {
          inset: 24px;
          animation-delay: 0.7s;
          border-color: rgba(236, 72, 153, 0.14);
        }

        .ontology-markdown {
          font-family: "DM Sans", sans-serif;
        }

        .ontology-markdown h2 {
          margin: 0 0 12px;
          font-family: "Syne", sans-serif;
          font-size: 15px;
          font-weight: 700;
          color: #e2e8f0;
        }

        .ontology-markdown h3 {
          margin: 18px 0 8px;
          font-size: 10px;
          letter-spacing: 0.22em;
          text-transform: uppercase;
          color: rgba(255, 255, 255, 0.35);
        }

        .ontology-markdown p,
        .ontology-markdown ol,
        .ontology-markdown ul {
          margin: 0 0 12px;
          font-size: 14px;
          line-height: 1.75;
          color: #e2e8f0;
        }

        .ontology-markdown ul {
          list-style: none;
          padding-left: 0;
        }

        .ontology-markdown ul li {
          position: relative;
          padding-left: 18px;
        }

        .ontology-markdown ul li::before {
          content: "◆";
          position: absolute;
          left: 0;
          top: 0;
          color: #a78bfa;
        }

        .ontology-markdown ol {
          padding-left: 18px;
        }

        .ontology-markdown table {
          width: 100%;
          border-collapse: collapse;
          overflow: hidden;
          border-radius: 12px;
          border: 1px solid rgba(255, 255, 255, 0.08);
        }

        .ontology-markdown th,
        .ontology-markdown td {
          padding: 10px 12px;
          text-align: left;
          font-size: 13px;
          border-bottom: 1px solid rgba(255, 255, 255, 0.06);
        }

        .ontology-markdown thead {
          background: rgba(255, 255, 255, 0.04);
        }

        .ontology-markdown tbody tr:nth-child(even) {
          background: rgba(139, 92, 246, 0.05);
        }

        .ontology-markdown strong {
          color: #e2e8f0;
          font-weight: 700;
        }

        .ontology-markdown em {
          color: rgba(255, 255, 255, 0.35);
          font-style: italic;
        }

        .ontology-markdown code {
          border-radius: 8px;
          background: rgba(255, 255, 255, 0.05);
          padding: 2px 6px;
          font-family: "JetBrains Mono", monospace;
          font-size: 12px;
        }

        .ontology-markdown pre {
          margin: 14px 0;
          overflow-x: auto;
          border-radius: 12px;
          border: 1px solid rgba(255, 255, 255, 0.06);
          background: rgba(0, 0, 0, 0.35);
          padding: 14px;
        }

        .ontology-markdown pre code {
          background: transparent;
          padding: 0;
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

            <div className="mt-6 rounded-[14px] border border-white/8 bg-white/4 p-3">
              <div className="text-[10px] uppercase tracking-[0.22em] text-[rgba(255,255,255,0.3)]">
                Quick Run
              </div>
              <div className="mt-3 space-y-2">
                {samplePrompts.map((sample) => (
                  <button
                    key={sample}
                    type="button"
                    onClick={() => handleSuggestion(sample)}
                    className="w-full rounded-[10px] border border-white/8 bg-[rgba(255,255,255,0.03)] px-3 py-2 text-left text-xs text-[rgba(255,255,255,0.7)] transition hover:border-[rgba(167,139,250,0.24)] hover:bg-[rgba(167,139,250,0.08)] hover:text-[#e2e8f0]"
                  >
                    {sample}
                  </button>
                ))}
              </div>
            </div>

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
                    Ontology Agent
                  </div>
                  <div className="mt-2 text-sm text-[rgba(255,255,255,0.35)]">
                    Premium chatbot for ontology queries and summaries.
                  </div>
                </div>
                <div className="rounded-full border border-white/8 bg-white/4 px-3 py-1.5 text-[11px] uppercase tracking-[0.22em] text-[rgba(255,255,255,0.35)]">
                  {isStreaming ? "processing live" : "idle"}
                </div>
              </div>
            </div>

            <div className="mt-4 flex min-h-0 flex-1 flex-col gap-4">
              <section className="glass-panel flex min-h-[420px] flex-col overflow-hidden rounded-[14px]">
                <div className="border-b border-white/8 px-5 py-4">
                  <div className="flex items-center justify-between gap-3">
                    <div className="text-[10px] uppercase tracking-[0.24em] text-[rgba(255,255,255,0.35)]">
                      Result
                    </div>
                    <div className="font-['JetBrains_Mono',monospace] text-[10px] text-[rgba(255,255,255,0.2)]">
                      {totalElapsedMs > 0
                        ? `${(totalElapsedMs / 1000).toFixed(2)}s`
                        : "--"}
                    </div>
                  </div>
                </div>

                <div
                  ref={chatScrollRef}
                  className="flex-1 overflow-y-auto px-5 py-5"
                >
                  {runs.length === 0
                    ? <RingsEmptyState />
                    : runs.map((run) => renderRunCard(run))}
                </div>
              </section>

              <section className="glass-panel rounded-[14px] p-3">
                <div className="flex flex-col gap-3 md:flex-row md:items-end">
                  <div className="flex-1">
                    <textarea
                      value={prompt}
                      onChange={(event) => setPrompt(event.target.value)}
                      onKeyDown={handlePromptKeyDown}
                      placeholder="Ask about your ontology..."
                      disabled={isStreaming}
                      rows={3}
                      className="w-full resize-none rounded-[10px] border border-[rgba(255,255,255,0.08)] bg-[rgba(255,255,255,0.03)] px-4 py-3 text-sm text-[#e2e8f0] outline-none transition placeholder:text-[rgba(255,255,255,0.35)] focus:border-[rgba(167,139,250,0.4)] focus:shadow-[0_0_0_3px_rgba(139,92,246,0.1)] disabled:opacity-50"
                    />
                  </div>

                  <button
                    type="button"
                    onClick={handleSubmit}
                    disabled={isStreaming || !prompt.trim()}
                    className="inline-flex h-11 w-11 items-center justify-center rounded-[8px] bg-[linear-gradient(135deg,#6d28d9,#0e7490)] text-white transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-80"
                  >
                    {isStreaming
                      ? <Lock className="h-4 w-4" />
                      : <Send className="h-4 w-4" />}
                  </button>
                </div>
              </section>

              <section className="glass-panel min-h-[260px] rounded-[14px] p-4">
                <div className="flex flex-col gap-3 border-b border-white/8 pb-4 md:flex-row md:items-center md:justify-between">
                  <div className="flex items-center gap-2">
                    <span
                      className={cn(
                        "h-2.5 w-2.5 rounded-full",
                        activeRun?.response
                          ? "bg-[#6ee7b7]"
                          : isStreaming
                            ? "bg-amber-300 animate-pulse"
                            : "bg-[rgba(255,255,255,0.2)]",
                      )}
                    />
                    <span className="text-[10px] uppercase tracking-[0.24em] text-[rgba(255,255,255,0.35)]">
                      Processing Steps
                    </span>
                  </div>

                  {completionLabel
                    ? <div className="rounded-full border border-[rgba(16,185,129,0.2)] bg-[rgba(16,185,129,0.08)] px-3 py-1.5 text-xs text-[#6ee7b7]">
                        ✓ {completionLabel}
                      </div>
                    : null}
                </div>

                <div className="mt-4 space-y-4">
                  {steps.length === 0
                    ? <div className="rounded-[12px] border border-dashed border-white/8 bg-white/3 px-4 py-8 text-center text-sm text-[rgba(255,255,255,0.35)]">
                        Streaming reasoning steps from `/ontology/summary` will
                        appear here.
                      </div>
                    : steps.map((step, index) => {
                        const state =
                          step.status === "active"
                            ? "active"
                            : step.status === "done"
                              ? "done"
                              : "pending";

                        const elapsedMs =
                          step.elapsedMs ??
                          (step.status === "active"
                            ? now - step.startedAt
                            : null);

                        return (
                          <StepCard
                            key={step.id}
                            step={{ ...step, elapsedMs }}
                            index={index}
                            state={state}
                            onToggle={() =>
                              setSteps((current) =>
                                current.map((entry) =>
                                  entry.id === step.id
                                    ? { ...entry, expanded: !entry.expanded }
                                    : entry,
                                ),
                              )
                            }
                          />
                        );
                      })}
                </div>
              </section>
            </div>
          </main>
        </div>
      </div>
    </>
  );
}
