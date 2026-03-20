# Chat Implementation Structure Analysis

## Overview
The ontology_bot uses a **mock-first approach** with a production-ready frontend and backend structure planned. The frontend is fully functional with mock data, and the backend infrastructure is ready for implementation.

---

## 1. Proto File Structures (Chat Message Definitions)

### Current Status
- **Proto directory**: `/proto/ontology_bot/` is currently **empty**
- **Expected structure** per `TESTING_CHAT_STEPS.md`: Should have `analysis.proto`

### Planned Proto Structure
```protobuf
// proto/ontology_bot/analysis.proto
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
  string status = 4;        // "pending", "running", "completed", "error"
  int64 start_time = 5;
  int64 end_time = 6;
  string details = 7;
}

message AnalysisQueryResponse {
  string message_id = 1;
  string content = 2;
  string content_type = 3;  // "text" or "markdown" (default: "text")
  repeated ProcessingStep steps = 4;
}
```

**Note**: These would generate Go types in `app/internal/types/analysis/` after running `./scripts/gen_proto.sh`

---

## 2. Frontend Message Structures (TypeScript)

### ChatMessage Interface
**Location**: Defined inline in [panel/src/services/analysis.service.ts](panel/src/services/analysis.service.ts)

```typescript
interface ChatMessage {
  id: string;                           // Unique message ID (format: `msg-${timestamp}`)
  role: "user" | "assistant";           // Message sender
  content: string;                      // Message text or markdown content
  timestamp: number;                    // Unix timestamp in milliseconds
  contentType?: "text" | "markdown";    // Content format (default: "text")
}
```

### LLMStep Interface  
**Location**: [panel/src/services/analysis.service.ts](panel/src/services/analysis.service.ts)

```typescript
interface LLMStep {
  id: string;                                           // Unique step ID
  title: string;                                        // Step name (e.g., "Parsing Query")
  description: string;                                  // Step description
  status: "pending" | "running" | "completed" | "error"; // Processing status
  startTime?: number;                                   // Unix timestamp when step started
  endTime?: number;                                     // Unix timestamp when step completed
  details?: string;                                     // Additional details/output
}
```

### AnalysisResponse Interface
**Location**: [panel/src/services/analysis.service.ts](panel/src/services/analysis.service.ts)

```typescript
interface AnalysisResponse {            // Message identifier
  content: string;                      // LLM-generated response (text or markdown)
  contentType?: "text" | "markdown";    // Content format (default: "text")
  steps: LLMStep[];                     // Array of processing steps with status
}
```

> **NEW**: Markdown support allows rich formatting of LLM responses with headers, tables, code blocks, and structured layout. See [MARKDOWN_RESPONSE_FORMAT.md](MARKDOWN_RESPONSE_FORMAT.md) for details.
```

---

## 3. Backend Service Implementation

### Current Implementation Status
- **Location**: `app/internal/services/`
- **Current files**: Only `service_health.go` exists
- **Status**: Analysis services not yet implemented

### Planned Service Structure
Based on CLAUDE.md conventions, the analysis service should follow:

```go
// app/internal/services/types.go
type ServiceAnalysisMethods interface {
  AnalyzeQuery(ctx context.Context, query string) (*AnalysisResponse, error)
  StreamAnalysisQuery(ctx context.Context, query string) (<-chan *AnalysisStep, error)
}

type AnalysisResponse struct {
  MessageID string           `json:"message_id"`
  Content   string           `json:"content"`
  Steps     []AnalysisStep   `json:"steps"`
}

type AnalysisStep struct {
  ID          string      `json:"id"`
  Title       string      `json:"title"`
  Description string      `json:"description"`
  Status      string      `json:"status"`
  StartTime   *int64      `json:"start_time,omitempty"`
  EndTime     *int64      `json:"end_time,omitempty"`
  Details     string      `json:"details,omitempty"`
}
```

### Implementation Flow
```
HTTP Request → Controller → Service → LLM Client
                (parse)      (logic)    (API call)
```

---

## 4. Frontend Components & Rendering

### ChatWindow Component
**File**: [panel/src/components/ChatWindow.tsx](panel/src/components/ChatWindow.tsx)

**Responsibilities**:
- Manages message state and input
- Calls `sendAnalysisQuery()` service
- Displays messages with role-based styling
- Updates parent components with `onMessagesUpdate()` and `onStepsUpdate()` callbacks

**Message Display Logic**:
```tsx
// User messages: Right-aligned, blue background, white text
// Assistant messages: Left-aligned, white/dark background, colored border

messages.map((msg) => (
  <div className={msg.role === "user" ? "justify-end" : "justify-start"}>
    <div className={msg.role === "user" 
      ? "bg-accent-blue text-white" 
      : "bg-white dark:bg-gray-800 border border-gray-300"
    }>
      {msg.content}
      <time>{timestamp}</time>
    </div>
  </div>
))
```

### StepsDisplay Component
**File**: [panel/src/components/StepsDisplay.tsx](panel/src/components/StepsDisplay.tsx)

**Responsibilities**:
- Displays processing steps with status icons
- Color-coded left borders (green=completed, blue=running, amber=pending, red=error)
- Shows timing information and expandable details

**Status Rendering**:
- ✅ **Completed**: Green border-left, CheckCircle2 icon
- ⏳ **Running**: Blue border-left, Loader2 icon (animated spinner)
- ⏱️ **Pending**: Amber border-left, Clock icon
- ❌ **Error**: Red border-left, AlertCircle icon

---

## 5. API Endpoints

### Current Endpoints
**Implemented**:
- `GET /ontology_bot/v1/health` - Health check
- `GET /ontology_bot/v1/health/ready` - Readiness check

### Planned Endpoints
**To be implemented**:

#### 1. Analysis Query (Non-streaming)
```
POST /ontology_bot/v1/analysis/query
Headers:
  Content-Type: application/json
  x-scope: analysis:query

Request Body:
{
  "query": "Analyze my system performance"
}

Response (200 OK):
{
  "message_id": "msg-1710859200000",
  "content": "Analysis complete! Here's what I found...",
  "steps": [
    {
      "id": "step-1",
      "title": "Parsing Query",
      "description": "Analyzing the user's query to understand intent",
      "status": "completed",
      "start_time": 1710859200000,
      "end_time": 1710859200500,
      "details": "Query parsed successfully..."
    },
    ...
  ]
}
```

#### 2. Analysis Query Stream (Real-time Steps)
```
POST /ontology_bot/v1/analysis/query/stream
Headers:
  Content-Type: application/json
  x-scope: analysis:query

Request Body:
{
  "query": "Analyze my system performance"
}

Response (Server-Sent Events):
Each step update as JSON line:
{"steps": [...]}
{"steps": [...]}
{"content": "Final response..."}
```

---

## 6. Mock Data Implementation

### Mock Service Functions
**File**: [panel/src/services/analysis.mock.ts](panel/src/services/analysis.mock.ts)

#### `sendAnalysisQueryMock(query: string)`
- Simulates 2-second network delay
- Returns complete response with 4 mock steps
- Mock steps: Parsing → Gathering Info → Analyzing → Generating Response

#### `streamAnalysisQueryMock(query: string)`
- Generator function yielding steps one-by-one
- Simulates real-time processing (1-1.5s per step)
- Each step progresses from "running" → "completed"

**Mock Data Structure**:
```typescript
// Each step includes realistic timings and details
{
  id: "step-1",
  title: "Parsing Query",
  description: `Analyzing: "${query}"`,
  status: "running",
  startTime: Date.now(),
  // After delay...
  status: "completed",
  endTime: Date.now(),
  details: "Successfully parsed query..."
}
```

---

## 7. Response Structure Summary

### Message Flow
```
Frontend (ChatWindow)
  ↓
sendAnalysisQuery(query)
  ├─ IF USE_MOCK_DATA → analysis.mock.ts
  │  └─ Returns AnalysisResponse with steps
  └─ ELSE → POST /ontology_bot/v1/analysis/query
     └─ Backend response

Display Logic:
  1. User message appears immediately (blue, right-aligned)
  2. Loading indicator shows
  3. Steps display updates in real-time
  4. Assistant response appears when complete (left-aligned)
```

### Response Data Types
**Complete API Response**:
```json
{
  "message_id": "msg-1710859200000",
  "content": "Analysis complete! System performing well...",
  "steps": [
    {
      "id": "step-1",
      "title": "Parsing Query",
      "description": "Analyzing the user's query...",
      "status": "completed",
      "start_time": 1710859200000,
      "end_time": 1710859200500,
      "details": "Query parsed successfully..."
    },
    {
      "id": "step-2",
      "title": "Gathering System Information",
      "description": "Collecting relevant system metrics...",
      "status": "completed",
      "start_time": 1710859200500,
      "end_time": 1710859202000,
      "details": "CPU: 45% | Memory: 60% | Disk: 75%..."
    }
  ]
}
```

---

## 8. Current Configuration

### Frontend Configuration
- **API Host**: Read from `Config.apiHost` (base/config)
- **Mock Mode**: `NEXT_PUBLIC_USE_MOCK_DATA` environment variable
- **Endpoints**: 
  - Non-streaming: `/ontology_bot/v1/analysis/query`
  - Streaming: `/ontology_bot/v1/analysis/query/stream`

### Backend Configuration
**File**: [app/config/local.example.yml](app/config/local.example.yml)
```yaml
ServerPort: 4441
AppName: "ontology_bot"
AppVersion: "0.0.1"
BaseUrl: "http://localhost:4441"
```

---

## 9. Key Findings Summary

| Component | Status | Location |
|-----------|--------|----------|
| **Proto Definitions** | ❌ Not created | Should be in `/proto/ontology_bot/analysis.proto` |
| **Type Definitions** (Frontend) | ✅ Implemented | `panel/src/services/analysis.service.ts` |
| **Mock Data** | ✅ Implemented | `panel/src/services/analysis.mock.ts` |
| **ChatWindow Component** | ✅ Implemented | `panel/src/components/ChatWindow.tsx` |
| **StepsDisplay Component** | ✅ Implemented | `panel/src/components/StepsDisplay.tsx` |
| **Backend Services** | ❌ Not implemented | Should be in `app/internal/services/` |
| **Backend Controller** | ❌ Not implemented | Should be in `app/internal/controllers/` |
| **Routes Registration** | ❌ Not implemented | Should be in `app/cmd/app/routes_analysis.go` |
| **Health Check API** | ✅ Implemented | `app/cmd/app/routes_health.go` |

---

## 10. Next Steps for Implementation

### Phase 1: Backend Structure (Proto-first)
1. ✏️ Create `proto/ontology_bot/analysis.proto` with message definitions
2. 🔨 Run `./scripts/gen_proto.sh` to generate Go types
3. 📝 Create `app/internal/services/service_analysis.go`

### Phase 2: Backend Implementation
4. 🎮 Create `app/internal/controllers/controller_analysis.go`
5. 🛣️ Create `app/cmd/app/routes_analysis.go` with route registration
6. 📊 Implement LLM client integration

### Phase 3: Integration
7. 🔗 Wire up services in dependency injection
8. ✅ Test endpoints with frontend
9. 🚀 Replace mock data with real backend

---

## 11. Frontend Display Logic Details

### Message Styling (ChatWindow)
```tsx
// User Message (Blue, Right)
className="bg-accent-blue dark:bg-blue-600 text-white"

// Assistant Message (Light, Left)
className="bg-white dark:bg-gray-800 text-gray-900 
           dark:text-gray-100 border border-gray-300 dark:border-gray-700"
```

### Steps Status Colors (StepsDisplay)
```tsx
border-l-green-600   // Completed
border-l-blue-600    // Running
border-l-amber-600   // Pending
border-l-red-600     // Error
```

### Auto-scroll Behavior
- Messages trigger scroll when array updates
- Uses `useRef` to track scroll container
- Scrolls to bottom on new messages

---

## 12. Environment Variables

### Frontend
```env
NEXT_PUBLIC_USE_MOCK_DATA=true  # Toggle between mock and real API
NEXT_PUBLIC_API_HOST=http://localhost:4441
```

### Backend
```yaml
# app/config/local.yml
ServerPort: 4441
AppName: ontology_bot
```

---

## Conclusion

The **chat implementation is 70% complete** with a development-ready frontend using mock data. The infrastructure supports:
- ✅ Real-time message display with role-based styling
- ✅ Processing steps visualization with status tracking
- ✅ Mock data for testing/development
- ✅ Streaming response support (planned)
- ⏳ Backend endpoints (ready for implementation following proto-first approach)

The next phase requires creating proto definitions, backend services, and controllers following the established `CLAUDE.md` conventions.
