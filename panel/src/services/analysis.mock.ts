import type { AnalysisResponse, LLMStep } from "./analysis.service";

// Mock data simulating agent thinking steps
const mockSteps: LLMStep[] = [
  {
    id: "step-1",
    title: "Parsing Query",
    description: "Analyzing the user's query to understand intent",
    status: "completed",
    startTime: Date.now() - 5000,
    endTime: Date.now() - 4500,
    details: "Query parsed successfully. Intent: user is asking about system performance.",
  },
  {
    id: "step-2",
    title: "Gathering System Information",
    description: "Collecting relevant system metrics and data",
    status: "completed",
    startTime: Date.now() - 4500,
    endTime: Date.now() - 3000,
    details: "CPU: 45% | Memory: 60% | Disk: 75% | Network: Stable",
  },
  {
    id: "step-3",
    title: "Analyzing Patterns",
    description: "Running analysis on collected data",
    status: "completed",
    startTime: Date.now() - 3000,
    endTime: Date.now() - 1500,
    details: "Identified 3 potential bottlenecks in the current configuration.",
  },
  {
    id: "step-4",
    title: "Generating Recommendations",
    description: "Formulating actionable insights and recommendations",
    status: "completed",
    startTime: Date.now() - 1500,
    endTime: Date.now(),
    details: "Generated 5 optimization recommendations with priority scores.",
  },
];

export async function sendAnalysisQueryMock(
  query: string,
): Promise<AnalysisResponse> {
  // Simulate network delay
  await new Promise((resolve) => setTimeout(resolve, 2000));

  return {
    message_id: `msg-${Date.now()}`,
    content: `# System Analysis Report

## Summary
Analysis complete! Here's a comprehensive breakdown of your system performance for: **"${query}"**

---

## Key Findings

### 1. Performance Metrics
- **CPU Usage**: 45% (healthy)
- **Memory Allocation**: 60% (16GB/32GB)
- **Disk I/O**: Normal operations
- **Network Latency**: 5ms

### 2. System Health
✅ **Overall Status**: Healthy  
✅ **Uptime**: 99.8% last 30 days  
⚠️ **Warning**: Database queries could be optimized

---

## Recommendations

| Priority | Action | Expected Impact | Effort |
|----------|--------|-----------------|--------|
| High | Implement request caching | 35% improvement | 2 days |
| Medium | Optimize database indexes | 15% improvement | 1 day |
| Low | Update dependencies | 5% improvement | 3 days |

---

## Code Examples

### Recommended Cache Implementation
\`\`\`typescript
// Example: Add caching layer
const cache = new Map<string, CacheEntry>();

function getCachedData(key: string) {
  const entry = cache.get(key);
  if (entry && !entry.isExpired()) {
    return entry.data;
  }
  return null;
}
\`\`\`

---

## Next Steps

1. Review the detailed analysis steps below
2. Prioritize high-impact recommendations
3. Schedule performance optimization sprint
4. Monitor improvement metrics post-deployment

> **Note**: This analysis is based on current system state. Re-run analysis after implementing changes for updated metrics.`,
    contentType: "markdown",
    steps: mockSteps,
  };
}

// Streaming version (simulates real-time steps appearing one by one)
export async function* streamAnalysisQueryMock(query: string) {
  const allSteps: LLMStep[] = [];

  // Step 1 - Parsing
  console.log("Step 1: Starting parse");
  const step1: LLMStep = {
    id: "step-1",
    title: "Parsing Query",
    description: `Analyzing: "${query}"`,
    status: "running",
    startTime: Date.now(),
  };
  allSteps.push(step1);
  yield { steps: [...allSteps] };

  await new Promise((resolve) => setTimeout(resolve, 1000));
  step1.status = "completed";
  step1.endTime = Date.now();
  step1.details = `Successfully parsed query. Keywords found: ${query.split(" ").slice(0, 3).join(", ")}`;
  yield { steps: [...allSteps] };

  // Step 2 - Gathering Info
  console.log("Step 2: Gathering info");
  const step2: LLMStep = {
    id: "step-2",
    title: "Gathering System Information",
    description: "Collecting metrics from system...",
    status: "running",
    startTime: Date.now(),
  };
  allSteps.push(step2);
  yield { steps: [...allSteps] };

  await new Promise((resolve) => setTimeout(resolve, 1200));
  step2.status = "completed";
  step2.endTime = Date.now();
  step2.details = "System Metrics:\n- CPU Usage: 45%\n- Memory: 60% (16GB/32GB)\n- Disk I/O: Normal\n- Network Latency: 5ms";
  yield { steps: [...allSteps] };

  // Step 3 - Analyzing
  console.log("Step 3: Analyzing patterns");
  const step3: LLMStep = {
    id: "step-3",
    title: "Analyzing Patterns",
    description: "Processing collected data...",
    status: "running",
    startTime: Date.now(),
  };
  allSteps.push(step3);
  yield { steps: [...allSteps] };

  await new Promise((resolve) => setTimeout(resolve, 1500));
  step3.status = "completed";
  step3.endTime = Date.now();
  step3.details = "Analysis Results:\n- 3 Performance bottlenecks identified\n- 2 Security recommendations\n- 1 Optimization opportunity for database queries";
  yield { steps: [...allSteps] };

  // Step 4 - Generating Response
  console.log("Step 4: Generating response");
  const step4: LLMStep = {
    id: "step-4",
    title: "Generating Response",
    description: "Creating actionable recommendations...",
    status: "running",
    startTime: Date.now(),
  };
  allSteps.push(step4);
  yield { steps: [...allSteps] };

  await new Promise((resolve) => setTimeout(resolve, 800));
  step4.status = "completed";
  step4.endTime = Date.now();
  step4.details = "Generated 5 recommendations with priority scoring and implementation guides.";
  yield { steps: [...allSteps] };

  // Final message
  console.log("Analysis complete");
  yield {
    content: `Analysis Complete for "${query}"\n\nKey Findings:\n1. System is performing well overall\n2. Three optimization opportunities identified\n3. Implement caching layer for 35% improvement\n\nSee the agent steps on the left to review the complete analysis workflow.`,
  };
}
