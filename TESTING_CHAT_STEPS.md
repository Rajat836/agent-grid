# Testing Chat Section with Agent Steps

This guide explains how to test the chat functionality with agent steps display.

## Current Status

✅ **Frontend**: Fully implemented with:
- Chat window for user inputs
- Steps display component
- Proper state management
- Integration ready

❌ **Backend**: Endpoint not yet implemented
- Endpoint: `POST /ontology_bot/v1/analysis/query`
- Status: Needs implementation

---

## Option 1: Quick Testing with Mock Data (Fastest 🚀)

This approach lets you test immediately without needing a backend.

### Step 1: Create a Mock Service

Create a new file: `panel/src/services/analysis.mock.ts`

```typescript
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
    content: `Analysis complete! Here's what I found:\n\n1. **Performance Issue**: Your system is running at 60% memory capacity\n2. **Optimization**: Consider implementing a caching layer\n3. **Recommendation**: Upgrade RAM to 32GB for better performance\n4. **Estimated Impact**: ~35% performance improvement\n\nThe detailed steps of my analysis are shown on the right.`,
    steps: mockSteps,
  };
}

// Streaming version (simulates real-time steps)
export async function* streamAnalysisQueryMock(query: string) {
  const stepsToYield: LLMStep[] = [
    {
      id: "step-1",
      title: "Parsing Query",
      description: "Analyzing intent...",
      status: "running",
    },
  ];

  // Simulate step 1 running
  await new Promise((resolve) => setTimeout(resolve, 800));
  stepsToYield[0].status = "completed";
  stepsToYield[0].startTime = Date.now() - 1000;
  stepsToYield[0].endTime = Date.now();
  yield { steps: [stepsToYield[0]] };

  // Simulate step 2 running
  await new Promise((resolve) => setTimeout(resolve, 1000));
  const step2: LLMStep = {
    id: "step-2",
    title: "Gathering System Information",
    description: "Collecting metrics...",
    status: "completed",
    startTime: Date.now() - 1000,
    endTime: Date.now(),
  };
  yield { steps: [...stepsToYield, step2] };

  // Simulate step 3
  await new Promise((resolve) => setTimeout(resolve, 1200));
  const step3: LLMStep = {
    id: "step-3",
    title: "Analyzing Patterns",
    description: "Running analysis...",
    status: "completed",
    startTime: Date.now() - 1200,
    endTime: Date.now(),
  };
  yield { steps: [...stepsToYield, step2, step3] };

  // Simulate final step
  await new Promise((resolve) => setTimeout(resolve, 1000));
  const step4: LLMStep = {
    id: "step-4",
    title: "Generating Recommendations",
    description: "Creating insights...",
    status: "completed",
    startTime: Date.now() - 1000,
    endTime: Date.now(),
  };
  yield { steps: [...stepsToYield, step2, step3, step4] };

  // Final message
  yield {
    content:
      "Analysis complete! The steps show the analysis workflow. This demonstrated how agent steps are displayed in real-time.",
  };
}
```

### Step 2: Update the Analysis Service to Use Mock Data

Edit: `panel/src/services/analysis.service.ts`

Add this at the top of the file (after imports):

```typescript
// Enable this flag to test with mock data
const USE_MOCK_DATA = process.env.NEXT_PUBLIC_USE_MOCK_DATA === "true";
```

Then update the functions:

```typescript
export async function sendAnalysisQuery(
  query: string,
): Promise<AnalysisResponse> {
  if (USE_MOCK_DATA) {
    const { sendAnalysisQueryMock } = await import("./analysis.mock");
    return sendAnalysisQueryMock(query);
  }

  // ... existing code
}

export async function* streamAnalysisQuery(query: string) {
  if (USE_MOCK_DATA) {
    const { streamAnalysisQueryMock } = await import("./analysis.mock");
    yield* streamAnalysisQueryMock(query);
    return;
  }

  // ... existing code
}
```

### Step 3: Enable Mock Mode

Add to `.env.local`:

```
NEXT_PUBLIC_USE_MOCK_DATA=true
```

### Step 4: Test It

1. Run the frontend:
```bash
cd panel
npm run dev
```

2. Open http://localhost:3000/system-analysis in your browser

3. Type any question in the chat (e.g., "Analyze my system performance")

4. You should see:
   - Message appears in chat ✅
   - Steps display panel appears below ✅
   - Steps show completion with color coding ✅
   - Final response message appears ✅

---

## Option 2: Implement Backend Endpoint (For Production 🏗️)

If you want to properly test with a real backend, follow these steps:

### Step 1: Create Backend Proto Definition

Create: `proto/ontology_bot/analysis.proto`

```protobuf
syntax = "proto3";

package ontology_bot.analysis;

// POST: /ontology_bot/v1/analysis/query
message AnalysisQueryRequest {
  string query = 1;
}

message ProcessingStep {
  string id = 1;
  string title = 2;
  string description = 3;
  string status = 4; // "pending", "running", "completed", "error"
  int64 start_time = 5;
  int64 end_time = 6;
  string details = 7;
}

message AnalysisQueryResponse {
  string message_id = 1;
  string content = 2;
  repeated ProcessingStep steps = 3;
}
```

### Step 2: Generate Types

```bash
cd /home/abhi/workspace/fyscal/ontology_bot
./scripts/gen_proto.sh
```

### Step 3: Create Backend Handler

Create: `app/internal/controllers/controller_analysis.go`

```go
package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    typesanalysis "app/ontology/internal/types/analysis"
    "time"
    "fmt"
)

type ControllerAnalysisMethods interface {
    Query(c *gin.Context)
}

type controllerAnalysis struct{}

func NewControllerAnalysis() ControllerAnalysisMethods {
    return &controllerAnalysis{}
}

func (ca *controllerAnalysis) Query(c *gin.Context) {
    var req struct {
        Query string `json:"query" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Simulate processing steps
    steps := []typesanalysis.ProcessingStep{
        {
            Id:          "step-1",
            Title:       "Parsing Query",
            Description: "Analyzing your query...",
            Status:      "completed",
            StartTime:   time.Now().UnixMilli(),
            EndTime:     time.Now().UnixMilli() + 500,
            Details:     fmt.Sprintf("Successfully parsed query: '%s'", req.Query),
        },
        {
            Id:          "step-2",
            Title:       "Gathering System Metrics",
            Description: "Collecting relevant data...",
            Status:      "completed",
            StartTime:   time.Now().UnixMilli() + 500,
            EndTime:     time.Now().UnixMilli() + 1500,
            Details:     "Collected CPU, Memory, and Network metrics",
        },
        {
            Id:          "step-3",
            Title:       "Running Analysis",
            Description: "Processing data patterns...",
            Status:      "completed",
            StartTime:   time.Now().UnixMilli() + 1500,
            EndTime:     time.Now().UnixMilli() + 2500,
            Details:     "Identified 3 optimization opportunities",
        },
        {
            Id:          "step-4",
            Title:       "Generating Report",
            Description: "Creating recommendations...",
            Status:      "completed",
            StartTime:   time.Now().UnixMilli() + 2500,
            EndTime:     time.Now().UnixMilli() + 3000,
            Details:     "Generated analysis report with 5 recommendations",
        },
    }

    response := typesanalysis.AnalysisQueryResponse{
        MessageId: fmt.Sprintf("msg-%d", time.Now().UnixNano()),
        Content: fmt.Sprintf(`Analysis Results for: "%s"

Your system analysis shows:
1. Performance score: 78/100
2. Memory usage is optimal
3. CPU throttling detected on Core 3
4. Network latency: 5ms average

See the steps on the right for detailed breakdown.`, req.Query),
        Steps: steps,
    }

    c.JSON(http.StatusOK, response)
}
```

### Step 4: Register Route

Create: `app/cmd/app/routes_analysis.go`

```go
package main

import (
    "github.com/gin-gonic/gin"
    "app/internal/controllers"
)

const (
    ScopeAnalysisQuery = "analysis:query"
)

func registerAnalysisRoutes(router *gin.Engine) {
    analysisController := controllers.NewControllerAnalysis()

    // Public routes
    public := router.Group("/ontology_bot/v1/analysis")
    {
        // POST /ontology_bot/v1/analysis/query
        public.POST("/query", analysisController.Query)
    }
}
```

### Step 5: Add Route Registration

In `app/cmd/app/routes.go`, add:

```go
func setupRoutes(router *gin.Engine) {
    // ... existing routes ...
    registerAnalysisRoutes(router)
}
```

### Step 6: Test with Backend

1. Start the backend:
```bash
cd app
make run
```

2. Update `.env.local`:
```
NEXT_PUBLIC_USE_MOCK_DATA=false
NEXT_PUBLIC_API_HOST=http://localhost:4441
```

3. Start the frontend:
```bash
cd panel
npm run dev
```

4. Test in browser at `http://localhost:3000/system-analysis`

---

## Manual Testing Checklist

- [ ] Message appears in chat when sent
- [ ] Steps panel appears below chat
- [ ] Steps show with correct status icons
- [ ] Completed steps are green
- [ ] Running steps have spinning animation
- [ ] Error steps are red
- [ ] Pending steps are gray
- [ ] Timing information displays correctly
- [ ] Details section shows for each step
- [ ] Chat input is disabled while loading
- [ ] Send button shows spinner while loading
- [ ] Multiple queries can be sent sequentially
- [ ] Previous messages remain in chat
- [ ] Steps clear when new query is sent

---

## Recommended Testing Path

**For immediate testing:**
1. Use Option 1 (Mock Data) - takes 5 minutes ⚡
2. Test all features
3. Check visual appearance

**For production:**
1. Implement Option 2 (Backend Endpoint)
2. Test with real data
3. Deploy

