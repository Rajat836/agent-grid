'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';

const domains = [
  {
    name: "Features", icon: "hexagon",
    tagline: "Discover and analyze feature performance",
    actions: ["get_features", "get_feature_metrics", "get_feature_instances", "create_feature"],
    enables: ["Feature discovery", "Performance analysis", "Activity tracking"]
  },
  {
    name: "Entities", icon: "nodes",
    tagline: "Understand workflows, transitions and node health",
    actions: ["get_entities", "get_transitions", "get_entity_metrics", "get_entity_apis"],
    enables: ["Workflow mapping", "Transition analysis", "Node performance"]
  },
  {
    name: "Services", icon: "layers",
    tagline: "Map ownership, deployments and criticality",
    actions: ["get_services", "update_service", "get_deployments", "assign_team"],
    enables: ["Ownership mapping", "Deployment tracking", "Criticality analysis"]
  },
  {
    name: "Teams", icon: "people",
    tagline: "Track ownership and operational accountability",
    actions: ["get_teams", "map_team_to_feature"],
    enables: ["Ownership visibility", "Accountability"]
  },
  {
    name: "APIs", icon: "brackets",
    tagline: "Monitor API performance across services",
    actions: ["get_apis", "update_api_metadata", "get_api_metrics"],
    enables: ["Performance monitoring", "Internal vs external"]
  },
  {
    name: "KPIs", icon: "chart",
    tagline: "Link business metrics to features and services",
    actions: ["get_kpis", "get_kpi_relationships", "map_kpi_to_feature"],
    enables: ["Business impact", "Metric-driven insights"]
  }
];

const examples: Record<string, any[]> = {
  Features: [{
    query: "Show performance of Customer Onboarding feature",
    plan: ["get_features", "get_feature_metrics"],
    output: "## Feature Performance\n\n- **Success Rate**: 94.2%\n- **Failure Rate**: 5.8%\n- **Avg Latency**: 312ms (p95: 890ms)"
  }],
  Entities: [{
    query: "Show entities and transitions for onboarding",
    plan: ["get_features", "get_entities", "get_transitions"],
    output: "## Entity Flow\n\n- **Start Node**: InitOnboarding\n- **Transitions**: 4 paths\n- **Terminal Nodes**: Success, Timeout, Failed"
  }],
  Ownership: [{
    query: "Which team owns the onboarding feature?",
    plan: ["get_features", "get_teams"],
    output: "## Ownership\n\n- **Feature**: Customer Onboarding\n- **Team**: Platform Core\n- **Contact**: platform@company.com"
  }],
  APIs: [{
    query: "Show API latency for auth service",
    plan: ["get_services", "get_apis", "get_api_metrics"],
    output: "## API Metrics – Auth Service\n\n- **p50 Latency**: 45ms\n- **p95 Latency**: 210ms\n- **Error Rate**: 0.3%"
  }],
  KPIs: [{
    query: "Which KPIs impact onboarding feature?",
    plan: ["get_kpis", "get_kpi_relationships"],
    output: "## KPI Impact\n\n- **Activation Rate** → directly linked\n- **Drop-off Rate** → inversely linked\n- **Time-to-Activate** → threshold: &lt;48hrs"
  }]
};

const navItems = [
  { id: 'how-it-works', label: 'How It Works' },
  { id: 'domains', label: 'Supported Domains' },
  { id: 'examples', label: 'Query Examples' },
  { id: 'execution', label: 'Execution Model' }
];

const xPositions = [60, 255, 450, 645, 840];
const yPositions = [80, 120, 160, 200, 240, 280, 320];
const stepData = [
  { from: 0, to: 1, text: "Enter query" },
  { from: 1, to: 2, text: "Send user query" },
  { from: 2, to: 3, text: "Generate execution plan" },
  { from: 3, to: 2, text: "Return JSON plan" },
  { from: 2, to: 4, text: "Call API (based on action)" },
  { from: 4, to: 2, text: "Return response + store in context" },
  { from: 2, to: 1, text: "Stream update → Final markdown response" },
];

const getIcon = (name: string) => {
  const props = { fill: "none", stroke: "currentColor", strokeWidth: "2", strokeLinecap: "round" as const, strokeLinejoin: "round" as const, width: 20, height: 20 };
  switch (name) {
    case 'hexagon': return <svg {...props} viewBox="0 0 24 24"><polygon points="12 2 21 7 21 17 12 22 3 17 3 7" /></svg>;
    case 'nodes': return <svg {...props} viewBox="0 0 24 24"><circle cx="18" cy="5" r="3" /><circle cx="6" cy="12" r="3" /><circle cx="18" cy="19" r="3" /><line x1="8.5" y1="10.5" x2="15.5" y2="6.5" /><line x1="8.5" y1="13.5" x2="15.5" y2="17.5" /></svg>;
    case 'layers': return <svg {...props} viewBox="0 0 24 24"><polygon points="12 2 2 7 12 12 22 7 12 2" /><polyline points="2 12 12 17 22 12" /><polyline points="2 17 12 22 22 17" /></svg>;
    case 'people': return <svg {...props} viewBox="0 0 24 24"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" /><circle cx="9" cy="7" r="4" /><path d="M23 21v-2a4 4 0 0 0-3-3.87" /><path d="M16 3.13a4 4 0 0 1 0 7.75" /></svg>;
    case 'brackets': return <svg {...props} viewBox="0 0 24 24"><polyline points="16 18 22 12 16 6" /><polyline points="8 6 2 12 8 18" /></svg>;
    case 'chart': return <svg {...props} viewBox="0 0 24 24"><path d="M3 3v18h18" /><path d="M18 9l-5 5-4-4-5 5" /></svg>;
    default: return null;
  }
};

interface ExecutionStepProps {
  num: number;
  title: string;
  desc: string;
  children: React.ReactNode;
  reverse?: boolean;
}

const ExecutionStep = ({ num, title, desc, children, reverse = false }: ExecutionStepProps) => (
  <div className={`exec-step ${reverse ? 'reverse' : ''} glass-panel p-8 rounded-2xl`}>
    <div className="flex-1">
      <div className="flex items-center gap-3 mb-4">
        <div className="w-7 h-7 rounded-full bg-violet-500/20 text-violet-400 flex items-center justify-center font-bold font-['Syne']">{num}</div>
        <h3 className="m-0 font-['Syne'] text-slate-200 text-xl">{title}</h3>
      </div>
      <p className="m-0 text-slate-400 text-sm leading-relaxed">{desc}</p>
    </div>
    <div className="flex-1 w-full min-w-0">
      {children}
    </div>
  </div>
);

const renderMockMarkdown = (md: string) => {
  return md.split('\n').map((line, i) => {
    if (line.startsWith('## ')) {
      return <div key={i} className="font-['Syne'] text-base text-slate-200 mb-3 font-semibold">{line.replace('## ', '')}</div>;
    }
    if (line.startsWith('- ')) {
      const parts = line.replace('- ', '').split('**');
      return (
        <div key={i} className="flex items-start gap-2 mb-2 text-sm text-slate-300">
          <div className="w-1 h-1 rounded-full bg-cyan-400 mt-1.5 flex-shrink-0" />
          <span className="leading-snug">
            {parts.map((part, j) => j % 2 === 1 ? <strong key={j} className="text-white font-semibold">{part.replace(/&lt;/g, '<')}</strong> : part.replace(/&lt;/g, '<'))}
          </span>
        </div>
      );
    }
    return <div key={i} className="h-1" />;
  });
};

export default function DocumentationPage({ onRunQuery }: { onRunQuery?: (query: string) => void }) {
  const router = useRouter();
  const [activeStep, setActiveStep] = useState(0);
  const [isPlaying, setIsPlaying] = useState(true);
  const [expandedDomain, setExpandedDomain] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState('Features');
  const [activeSection, setActiveSection] = useState('how-it-works');

  useEffect(() => {
    let interval: NodeJS.Timeout;
    if (isPlaying && activeStep <= stepData.length) {
      interval = setInterval(() => {
        setActiveStep((prev) => {
          if (prev >= stepData.length) {
            setIsPlaying(false);
            return prev;
          }
          return prev + 1;
        });
      }, 700);
    }
    return () => clearInterval(interval);
  }, [isPlaying, activeStep]);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        const visible = entries.find(e => e.isIntersecting);
        if (visible) {
          setActiveSection(visible.target.id);
        }
      },
      { threshold: 0.3, rootMargin: '0px 0px -50% 0px' }
    );

    navItems.forEach(item => {
      const el = document.getElementById(item.id);
      if (el) observer.observe(el);
    });

    return () => observer.disconnect();
  }, []);

  return (
    <div className="doc-scroll flex flex-col flex-1 h-full text-slate-200 font-['DM_Sans'] overflow-y-auto">
      <style>
        {`
          @keyframes streamFadeIn {
            0% { opacity: 0; transform: translateY(5px); }
            100% { opacity: 1; transform: translateY(0); }
          }
          .stream-line-1 { animation: streamFadeIn 0.3s ease forwards; animation-delay: 0.2s; opacity: 0; }
          .stream-line-2 { animation: streamFadeIn 0.3s ease forwards; animation-delay: 1.2s; opacity: 0; }
          .stream-line-3 { animation: streamFadeIn 0.3s ease forwards; animation-delay: 2.2s; opacity: 0; }
          
          .exec-step { display: flex; align-items: center; gap: 32px; }
          .exec-step.reverse { flex-direction: row-reverse; }
          
          @media (max-width: 900px) {
            .nav-sidebar { display: none !important; }
          }
          @media (max-width: 768px) {
            .exec-step, .exec-step.reverse { flex-direction: column; align-items: flex-start; }
            .main-grid { grid-template-columns: 1fr !important; }
            .svg-container { overflow-x: auto; }
          }
          
          .doc-scroll::-webkit-scrollbar { width: 6px; height: 6px; }
          .doc-scroll::-webkit-scrollbar-track { background: rgba(0,0,0,0.1); }
          .doc-scroll::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.1); border-radius: 4px; }
        `}
      </style>

      <div className="flex p-12 gap-12 relative">
        <div className="flex-1 flex flex-col gap-20 max-w-4xl">
          
          {/* SECTION 1: How It Works */}
          <section id="how-it-works" className="scroll-mt-10">
            <h2 className="font-['Syne'] text-2xl text-white mb-6">How It Works</h2>
            
            <div className="glass-panel rounded-2xl p-6 overflow-hidden">
              <div className="flex justify-end mb-4">
                <button 
                  onClick={() => {
                    if (activeStep >= stepData.length) {
                      setActiveStep(0);
                      setIsPlaying(true);
                    } else {
                      setIsPlaying(!isPlaying);
                    }
                  }}
                  className="bg-violet-500/10 border border-violet-500/40 text-violet-400 px-4 py-1.5 rounded-md cursor-pointer font-['JetBrains_Mono'] text-xs"
                >
                  {activeStep >= stepData.length ? 'Replay' : isPlaying ? 'Pause' : 'Play'}
                </button>
              </div>
              
              <div className="svg-container doc-scroll w-full">
                <svg width="100%" height="400" viewBox="0 0 900 400" className="min-w-[800px]">
                  {xPositions.map((x, i) => (
                    <g key={`col-${i}`}>
                      <line x1={x} y1="30" x2={x} y2="350" stroke="rgba(255,255,255,0.1)" strokeDasharray="4 4" />
                      <text x={x} y="20" fill="#e2e8f0" fontSize="13" className="font-['Syne']" textAnchor="middle">{["User", "Frontend", "Go Backend", "LLM", "External APIs"][i]}</text>
                    </g>
                  ))}
                  
                  <rect x="420" y="140" width="450" height="170" fill="rgba(167,139,250,0.02)" stroke="rgba(167,139,250,0.3)" strokeWidth="1" strokeDasharray="4 4" rx="8" />
                  <text x="430" y="155" fill="#a78bfa" fontSize="11" className="font-['JetBrains_Mono']">loop: for each step in plan</text>
                  
                  {stepData.map((step, i) => {
                    const isVisible = activeStep >= i + 1;
                    const isActive = activeStep === i + 1;
                    const x1 = xPositions[step.from];
                    const x2 = xPositions[step.to];
                    const y = yPositions[i];
                    const direction = x2 > x1 ? 1 : -1;
                    const color = isActive ? '#a78bfa' : '#22d3ee';
                    
                    if (!isVisible) return null;
                    
                    const centerX = Math.min(x1, x2) + Math.abs(x2 - x1) / 2;
                    const width = step.text.length * 6.8;
                    
                    return (
                      <g key={`step-${i}`} className="transition-all duration-300 ease-in-out">
                        <line x1={x1} y1={y} x2={x2 - direction * 5} y2={y} stroke={color} strokeWidth="2" />
                        <polygon points={`${x2},${y} ${x2 - direction * 6},${y - 4} ${x2 - direction * 6},${y + 4}`} fill={color} />
                        <rect x={centerX - width / 2} y={y - 20} width={width} height="20" fill="#0a0a18" stroke={color} strokeWidth="1" rx="10" />
                        <text x={centerX} y={y - 6} fill="#e2e8f0" fontSize="11" className="font-['DM_Sans']" textAnchor="middle">{step.text}</text>
                      </g>
                    );
                  })}
                </svg>
              </div>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
              <div className="glass-panel rounded-xl p-4">
                <div className="font-['Syne'] text-slate-200 mb-2 font-semibold">LLM decides</div>
                <div className="text-xs text-slate-400">which APIs, in what order, whether chaining is needed</div>
              </div>
              <div className="glass-panel rounded-xl p-4">
                <div className="font-['Syne'] text-slate-200 mb-2 font-semibold">Backend executes</div>
                <div className="text-xs text-slate-400">resolves mappings, calls APIs, stores context</div>
              </div>
              <div className="glass-panel rounded-xl p-4">
                <div className="font-['Syne'] text-slate-200 mb-2 font-semibold">You receive</div>
                <div className="text-xs text-slate-400">structured Markdown analysis, streamed live</div>
              </div>
            </div>
          </section>

          {/* SECTION 2: Supported Domains */}
          <section id="domains" className="scroll-mt-10">
            <h2 className="font-['Syne'] text-2xl text-white mb-6">Supported Domains</h2>
            <div className="main-grid grid grid-cols-2 gap-4">
              {domains.map((d, i) => {
                const isExpanded = expandedDomain === d.name;
                return (
                  <div key={i} className="glass-panel rounded-xl p-5 cursor-pointer transition-all duration-200 ease-in-out" onClick={() => setExpandedDomain(isExpanded ? null : d.name)}>
                    <div className="flex items-center gap-3">
                      <div className="text-violet-400">{getIcon(d.icon)}</div>
                      <div className="font-['Syne'] text-base font-semibold">{d.name}</div>
                    </div>
                    <div className="text-xs text-slate-400 my-2">{d.tagline}</div>
                    
                    <div className="flex flex-wrap gap-1.5 mb-3">
                      {d.enables.map(e => <span key={e} className="text-xs text-cyan-400 bg-cyan-500/10 px-2 py-0.5 rounded-full">{e}</span>)}
                    </div>

                    {isExpanded && (
                      <div className="mt-4 pt-4 border-t border-white/10">
                        <div className="text-xs mb-2 text-slate-300">Key Actions</div>
                        <div className="flex flex-wrap gap-2">
                          {d.actions.map(act => (
                            <span key={act} className="font-['JetBrains_Mono'] text-xs bg-violet-500/10 border border-violet-500/20 text-violet-400 px-2 py-1 rounded-md">
                              {act}
                            </span>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          </section>

          {/* SECTION 3: Query Examples */}
          <section id="examples" className="scroll-mt-10">
            <h2 className="font-['Syne'] text-2xl text-white mb-6">Query Examples</h2>
            <div className="glass-panel rounded-xl overflow-hidden">
              <div className="doc-scroll flex gap-2 p-4 border-b border-white/10 overflow-x-auto">
                {Object.keys(examples).map(tab => (
                  <button key={tab} onClick={() => setActiveTab(tab)} className={`border-solid border ${activeTab === tab ? 'bg-violet-500/15 border-violet-500/40' : 'border-transparent'} ${activeTab === tab ? 'text-violet-400' : 'text-slate-200'} px-4 py-1.5 rounded-full text-sm cursor-pointer font-['DM_Sans'] whitespace-nowrap`}>
                    {tab}
                  </button>
                ))}
              </div>
              <div className="p-6">
                {examples[activeTab].map((ex, i) => (
                  <div key={i} className={i > 0 ? 'mt-8' : ''}>
                    <div className="bg-cyan-500/10 border border-cyan-500/20 text-cyan-400 p-3 rounded-xl rounded-br-none inline-block mb-4">
                      {ex.query}
                    </div>
                    <div className="flex items-center gap-2 mb-4 flex-wrap">
                      {ex.plan.map((p: string, j: number) => (
                        <React.Fragment key={p}>
                          <span className="font-['JetBrains_Mono'] text-xs bg-white/5 px-2 py-1 rounded-md border border-white/10">{p}</span>
                          {j < ex.plan.length - 1 && <span className="text-violet-500/50">──▶</span>}
                        </React.Fragment>
                      ))}
                    </div>
                    <div className="bg-[#05050f] border border-white/5 rounded-lg p-4 mb-4">
                      {renderMockMarkdown(ex.output)}
                    </div>
                    <div className="flex justify-end">
                      <button onClick={() => {
                        if (onRunQuery) onRunQuery(ex.query);
                        else router.push('/system-analysis?q=' + encodeURIComponent(ex.query));
                      }} className="bg-transparent border-none text-violet-400 text-sm cursor-pointer py-2 px-0">
                        Try this query ↗
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </section>

          {/* SECTION 4: Execution Model */}
          <section id="execution" className="scroll-mt-10 pb-12">
            <h2 className="font-['Syne'] text-2xl text-white mb-6">Execution Model</h2>
            <div className="flex flex-col gap-8">
              
              <ExecutionStep num={1} title="LLM generates an execution plan" desc="The backend sends your query + all available actions to the LLM. The LLM returns a structured JSON plan deciding which APIs to call and in what order.">
                <div className="doc-scroll bg-[#05050f] p-4 rounded-lg border border-white/5 overflow-x-auto">
                  <pre className="m-0 font-['JetBrains_Mono'] text-xs text-slate-200">
{`{
  "plan": [
    { "step": 1, "action": "get_features" },
    { "step": 2, "action": "get_feature_metrics" }
  ]
}`}
                  </pre>
                </div>
              </ExecutionStep>

              <ExecutionStep num={2} title="Backend resolves input mappings" desc="Each step's output feeds the next. The backend automatically extracts IDs and parameters from prior results — no manual chaining needed." reverse>
                <div className="flex items-center gap-4 flex-wrap justify-center">
                  <div className="bg-violet-500/10 border border-violet-500/30 p-2 px-3 rounded-md text-violet-400 font-['JetBrains_Mono'] text-sm">
                    step1.result[0].id
                  </div>
                  <div className="text-white/40 text-xs flex items-center gap-2 font-['DM_Sans']">
                    <span>──▶</span>
                    <span className="px-2 py-1 bg-white/5 rounded-2xl">injected as</span>
                    <span>──▶</span>
                  </div>
                  <div className="bg-cyan-500/10 border border-cyan-500/30 p-2 px-3 rounded-md text-cyan-400 font-['JetBrains_Mono'] text-sm">
                    step2.feature_id
                  </div>
                </div>
              </ExecutionStep>

              <ExecutionStep num={3} title="Results stored, updates streamed" desc="Each API response is saved to a shared context object. The frontend receives live status updates as each step executes.">
                <div className="doc-scroll flex flex-col gap-4 bg-[#05050f] p-4 rounded-lg border border-white/5 overflow-x-auto">
                  <pre className="m-0 font-['JetBrains_Mono'] text-xs text-slate-200">
{`{
  "step1": { "id": "feat_001", "name": "Customer Onboarding" },
  "step2": { "success_rate": 0.942, "latency_p95": 890 }
}`}
                  </pre>
                  <div className="pt-4 border-t border-white/10 font-['JetBrains_Mono'] text-xs text-slate-400 flex flex-col gap-2">
                    <div className="stream-line-1 flex gap-2 items-center">
                      <span className="text-violet-400">▶</span> Fetching features...
                    </div>
                    <div className="stream-line-2 flex gap-2 items-center">
                      <span className="text-violet-400">▶</span> Getting metrics...
                    </div>
                    <div className="stream-line-3 flex gap-2 items-center">
                      <span className="text-cyan-400">✓</span> Analysis complete
                    </div>
                  </div>
                </div>
              </ExecutionStep>

            </div>
          </section>

        </div>

        {/* Sticky Nav Sidebar */}
        <div className="nav-sidebar w-52 flex-shrink-0">
          <div className="sticky top-6 flex flex-col gap-3">
            {navItems.map(item => (
              <div key={item.id} className="flex items-center gap-2.5 cursor-pointer" onClick={() => document.getElementById(item.id)?.scrollIntoView({ behavior: 'smooth' })}>
                <div className={`w-2 h-2 rounded-full border transition-all duration-200 ease-in-out ${activeSection === item.id ? 'bg-violet-400 border-violet-400' : 'border-white/20'}`} />
                <span className={`text-sm transition-all duration-200 ease-in-out ${activeSection === item.id ? 'text-violet-400 font-semibold' : 'text-slate-400'}`}>{item.label}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
