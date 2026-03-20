# Markdown Response Format Guide

## Overview
The LLM backend can now return responses formatted as **Markdown** for rich text rendering on the frontend. This enables better formatting, tables, code blocks, and structured content display.

---

## Response Format

### JSON Structure
```json
{
  "message_id": "msg-1710859200000",
  "content": "# Main Title\n\n## Section\n\nMarkdown content here...",
  "contentType": "markdown",
  "steps": [
    {
      "id": "step-1",
      "title": "Processing Step",
      "description": "What this step does",
      "status": "completed",
      "start_time": 1710859200000,
      "end_time": 1710859200500,
      "details": "Step details"
    }
  ]
}
```

### Content Type Field
- **`contentType`** (optional): 
  - `"markdown"` — Content will be rendered as Markdown
  - `"text"` or omitted — Plain text display

---

## Supported Markdown Elements

### 1. Headers
```markdown
# H1 Header
## H2 Header
### H3 Header
```

### 2. Text Formatting
```markdown
**bold text**
*italic text*
***bold italic***
~~strikethrough~~
```

### 3. Lists

**Unordered List:**
```markdown
- Item 1
- Item 2
  - Nested item
```

**Ordered List:**
```markdown
1. First item
2. Second item
3. Third item
```

### 4. Code Blocks

**Inline Code:**
```markdown
Use `const x = 10;` in your code
```

**Code Block (with language):****
````markdown
```typescript
interface User {
  id: string;
  name: string;
}
```
````

### 5. Tables (GitHub Flavored Markdown)
```markdown
| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |
| Cell 4   | Cell 5   | Cell 6   |
```

### 6. Blockquotes
```markdown
> This is a blockquote
> 
> Multiple paragraphs supported
```

### 7. Links
```markdown
[Link Text](https://example.com)
[Internal Link](./relative/path.md)
```

### 8. Horizontal Rules
```markdown
---
***
___
```

### 9. Task Lists (GitHub Flavored Markdown)
```markdown
- [x] Completed task
- [ ] Incomplete task
- [x] Another done task
```

---

## CSS Classes Applied to Elements

The `MarkdownRenderer` component applies Tailwind CSS classes:

| Element | Classes |
|---------|---------|
| `<h1>` | `text-lg font-bold my-2` |
| `<h2>` | `text-base font-bold my-2` |
| `<h3>` | `text-sm font-bold my-2` |
| `<ul>` | `list-disc list-inside space-y-1 my-2` |
| `<ol>` | `list-decimal list-inside space-y-1 my-2` |
| `<code>` (inline) | `bg-gray-200 dark:bg-gray-700 px-1.5 py-0.5 rounded font-mono text-xs` |
| `<code>` (block) | `block bg-gray-900 dark:bg-black text-gray-100 p-3 rounded-lg font-mono text-xs` |
| `<table>` | `border-collapse border border-gray-300 dark:border-gray-600 text-sm` |
| `<blockquote>` | `border-l-4 border-gray-400 dark:border-gray-600 pl-3 italic text-gray-600 dark:text-gray-400` |
| `<a>` | `text-blue-600 dark:text-blue-400 underline hover:opacity-80` |

---

## Backend Implementation Example

### Go Controller Response
```go
// POST /ontology_bot/v1/analysis/query
func (c ControllerAnalysisMethods) Query(ctx *gin.Context) {
  var req AnalysisQueryRequest
  if err := ctx.BindJSON(&req); err != nil {
    c.HandleError(ctx, err, "Invalid request", http.StatusBadRequest)
    return
  }

  // Call LLM service
  response, err := c.service.AnalyzeQuery(ctx.Request.Context(), req.Query)
  if err != nil {
    c.HandleError(ctx, err, "Analysis failed", http.StatusInternalServerError)
    return
  }

  // Return markdown response
  ctx.JSON(http.StatusOK, gin.H{
    "message_id": response.MessageID,
    "content": response.MarkdownContent,
    "contentType": "markdown",  // <-- Important!
    "steps": response.Steps,
  })
}
```

### Service Implementation
```go
type ServiceAnalysisMethods interface {
  AnalyzeQuery(ctx context.Context, query string) (*AnalysisResponse, error)
}

type AnalysisResponse struct {
  MessageID       string           `json:"message_id"`
  MarkdownContent string           `json:"content"`
  ContentType     string           `json:"contentType"`
  Steps           []ProcessingStep `json:"steps"`
}

func (s *service) AnalyzeQuery(ctx context.Context, query string) (*AnalysisResponse, error) {
  // Call LLM API (e.g., Claude, GPT-4)
  llmResponse := s.client.CallLLM(ctx, query)
  
  // LLM should return markdown formatted content
  return &AnalysisResponse{
    MessageID:       fmt.Sprintf("msg-%d", time.Now().UnixMilli()),
    MarkdownContent: llmResponse.MarkdownText,
    ContentType:     "markdown",
    Steps:           ProcessingSteps...,
  }, nil
}
```

---

## Frontend Implementation

### Using MarkdownRenderer
```typescript
import { MarkdownRenderer } from "@/components/MarkdownRenderer";

export function MyComponent() {
  const markdownContent = `
# Hello
This is **bold** text.
  `;

  return (
    <MarkdownRenderer 
      content={markdownContent}
      className="text-gray-900 dark:text-gray-100"
    />
  );
}
```

### In ChatWindow
- Messages with `contentType === "markdown"` are rendered with `<MarkdownRenderer>`
- Messages with `contentType === "text"` or no contentType are rendered as plain text
- User messages are always plain text

---

## Example Full Response

```json
{
  "message_id": "msg-1710859200000",
  "content": "# System Performance Analysis\n\n## Executive Summary\nYour system is performing **excellently** with the following metrics:\n\n### Performance Metrics\n\n| Metric | Value | Status |\n|--------|-------|--------|\n| CPU Usage | 45% | ✅ Optimal |\n| Memory | 16GB/32GB | ✅ Healthy |\n| Disk I/O | Normal | ✅ Good |\n\n## Recommendations\n\n1. **Implement Caching**\n   ```typescript\n   const cache = new Map();\n   ```\n2. **Optimize Queries**\n   - Add indexes\n   - Use pagination\n\n> **Note**: Monitor performance after changes\n",
  "contentType": "markdown",
  "steps": [
    {
      "id": "step-1",
      "title": "Parsing Query",
      "description": "Analyzing request",
      "status": "completed",
      "start_time": 1710859200000,
      "end_time": 1710859200500
    },
    {
      "id": "step-2",
      "title": "Analysis",
      "description": "Running analysis",
      "status": "completed",
      "start_time": 1710859200500,
      "end_time": 1710859202000
    }
  ]
}
```

---

## Dark Mode Support

The `MarkdownRenderer` automatically adapts to dark mode:
- Uses `dark:` Tailwind variants
- Code blocks use dark backgrounds in dark mode
- Links use appropriate colors for visibility
- Tables adjust borders and backgrounds

---

## Migration Path

### Step 1: Add contentType to Response
Start returning `"contentType": "markdown"` in your API responses

### Step 2: Frontend Recognition
The frontend automatically detects and renders markdown based on contentType

### Step 3: Gradual Adoption
- Initially only critical responses use markdown
- Gradually migrate other response types
- Keep backwards compatibility with plain text

---

## Best Practices

1. **Always set `contentType`**: Explicitly declare the format
2. **Structure content**: Use headers to organize information
3. **Use tables**: For comparisons and structured data
4. **Code examples**: Include when relevant
5. **Blockquotes**: Highlight important notes or warnings
6. **Validation**: Test markdown rendering on both light and dark themes

---

## Troubleshooting

### Markdown not rendering?
- Check `contentType: "markdown"` is set in response
- Verify markdown syntax is valid
- Check browser console for errors

### Styling issues?
- Ensure TailwindCSS is properly configured
- Check for CSS conflicts with other components
- Verify dark mode is enabled if needed

### Performance concerns?
- Large markdown content is fine (100KB+)
- Rendering is optimized by `react-markdown`
- Pagination not needed for response content

