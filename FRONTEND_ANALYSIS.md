# Frontend (Panel) Analysis - Ontology Bot

## Executive Summary

The frontend is a Next.js 13+ App Router application using shadcn/Radix UI components with Tailwind CSS. It follows a clean, modular structure with a **sidebar-based navigation layout**. There are **minimal overlap issues** currently, but the layout uses flexbox constraints that could be optimized.

---

## 1. Main Layout Structure

### Root Layout (`src/app/layout.tsx`)
- **Minimal root layout** - only sets up:
  - Geist font family (sans + mono variants)
  - Sonner toast notifications
  - Basic `<html>` and `<body>` structure

### Root Page (`src/app/page.tsx`)
- **Simple redirect** - routes all traffic to `/system-analysis`

### Layout Hierarchy
```
Root Layout (layout.tsx)
├── Root Page (page.tsx) → redirects to /system-analysis
├── System Analysis Route
│   ├── page.tsx (System Analysis Page)
│   ├── SidebarLayout (wrapper)
│   ├── ChatWindow (primary)
│   └── StepsDisplay (secondary)
└── Documentation Route
    ├── page.tsx (Documentation Page)
    ├── SidebarLayout (wrapper)
    └── Placeholder content
```

---

## 2. Component Structure & File Organization

### Key Components Location

**Main Components** (`src/components/`):
- `ChatWindow.tsx` - Chat interface for LLM analysis
- `StepsDisplay.tsx` - Processing steps visualization
- `app-sidebar.tsx` - Navigation sidebar wrapper
- `sidebar-layout.tsx` - Layout provider for sidebar + content
- `ui/` - 50+ shadcn/Radix UI components (buttons, dialogs, forms, etc.)

### Navigation Items (`src/lib/sidebar.ts`)
```typescript
navItems = [
  { title: "System Analysis", icon: BarChart3, url: "/system-analysis" },
  { title: "Documentation", icon: FileText, url: "/documentation" }
];
secondaryNavItems = []; // Empty
```

---

## 3. Layout Deep Dive - System Analysis Page

### Structure Visualization
```
┌─────────────────────────────────────────────────────────────┐
│                    SidebarLayout                             │
├──────────────┬──────────────────────────────────────────────┤
│   SIDEBAR    │           SidebarInset                        │
│ (left side)  │  ┌────────────────────────────────────────┐  │
│              │  │ System Analysis Page                   │  │
│ • Analysis   │  │ ┌──────────────────────────────────┐   │  │
│ • Docs       │  │ │ <h1>System Analysis</h1>        │   │  │
│              │  │ │                                  │   │  │
│              │  │ │ ┌─────────────────────────────┐ │   │  │
│              │  │ │ │  ChatWindow                 │ │   │  │
│              │  │ │ │  (flex-1, min-h-0)         │ │   │  │
│              │  │ │ │                             │ │   │  │
│              │  │ │ │ - Scrollable area           │ │   │  │
│              │  │ │ │ - Message list              │ │   │  │
│              │  │ │ │ - Input field               │ │   │  │
│              │  │ │ └─────────────────────────────┘ │   │  │
│              │  │ │                                  │   │  │
│              │  │ │ ┌─────────────────────────────┐ │   │  │
│              │  │ │ │ StepsDisplay (conditional)  │ │   │  │
│              │  │ │ │ (flex-1, min-h-0)          │ │   │  │
│              │  │ │ │                             │ │   │  │
│              │  │ │ │ - Shows when steps exist    │ │   │  │
│              │  │ │ │ - Loading animation         │ │   │  │
│              │  │ │ └─────────────────────────────┘ │   │  │
│              │  │ └──────────────────────────────────┘   │  │
│              │  └────────────────────────────────────────┘  │
└──────────────┴──────────────────────────────────────────────┘
```

### Flexbox Constraints (`src/app/system-analysis/page.tsx`)
```typescript
<SidebarLayout>
  <div className="flex flex-col h-full gap-4 px-6 py-4">
    <h1 className="text-2xl font-bold shrink-0">System Analysis</h1>
    
    <div className="flex flex-col min-h-0 gap-4 flex-1">
      {/* ChatWindow - Takes 50% or more space */}
      <div className="flex-1 min-h-0">
        <ChatWindow />
      </div>
      
      {/* StepsDisplay - Only shows when needed */}
      {(steps.length > 0 || isLoading) && (
        <div className="flex-1 min-h-0 border rounded-lg p-4 bg-card overflow-y-auto">
          <StepsDisplay steps={steps} isLoading={isLoading} />
        </div>
      )}
    </div>
  </div>
</SidebarLayout>
```

**Key Classes:**
- `h-full` - Full height of viewport
- `flex-1` - Take available space (equal distribution)
- `min-h-0` - Critical for scroll areas to work properly
- `gap-4` - 16px spacing between components
- `overflow-hidden` - Prevent overflow at container level
- `overflow-y-auto` - Allow scrolling within components

---

## 4. Overlap Analysis & Issues

### ✅ Current State: MINIMAL OVERLAP

**Reason:** The layout uses proper flexbox constraints with `min-h-0` to ensure scroll areas work correctly.

### Potential Issues & Recommendations

| Issue | Location | Severity | Recommendation |
|-------|----------|----------|-----------------|
| **Flex Space Competition** | ChatWindow + StepsDisplay | LOW | Both use `flex-1`, so they split space equally once both are visible. On large screens, may want different ratios (e.g., 60/40). |
| **No Header** | System Analysis Page | LOW | Page has title but no sticky header. Could benefit from fixed header with controls. |
| **Limited Mobile Layout** | Sidebar | MEDIUM | Sidebar collapses to icon-only on mobile, but chat window stays full-width. May need responsive adjustments. |
| **No Z-Index Management** | All components | LOW | No popovers/modals currently, but should document z-index strategy if added. |
| **Scroll Area Performance** | ChatWindow + StepsDisplay | LOW | Multiple nested scroll areas could cause jank if message volume grows large. |

### Recommended Flex Ratios for Future
```typescript
// Current: 50/50 split after both visible
// Recommended: 60% chat, 40% steps
<div className="flex-[0.6] min-h-0">
  <ChatWindow />
</div>
{showSteps && (
  <div className="flex-[0.4] min-h-0 border rounded-lg p-4 bg-card overflow-y-auto">
    <StepsDisplay steps={steps} />
  </div>
)}
```

---

## 5. Color Scheme & Styling

### Design System
- **Framework:** Tailwind CSS + shadcn/Radix UI (headless components)
- **Color System:** OKLCH color space (modern, perceptually uniform)
- **Fonts:** Geist (sans), Geist Mono
- **Rounded Corners:** `--radius: 0.625rem` (10px)

### Color Variables (`src/app/globals.css`)

#### Light Mode (`:root`)
```css
--background: oklch(1 0 0)                    /* White */
--foreground: oklch(0.145 0 0)                /* Dark gray/black */
--primary: oklch(0.205 0 0)                   /* Near black */
--primary-foreground: oklch(0.985 0 0)        /* Near white */
--secondary: oklch(0.97 0 0)                  /* Very light gray */
--secondary-foreground: oklch(0.205 0 0)      /* Dark */
--muted: oklch(0.97 0 0)                      /* Light gray (messages) */
--muted-foreground: oklch(0.556 0 0)          /* Medium gray */
--accent: oklch(0.97 0 0)                     /* Light gray */
--destructive: oklch(0.577 0.245 27.325)      /* Red */
```

#### Dark Mode (`.dark`)
```css
--background: oklch(0.145 0 0)                /* Dark gray/black */
--foreground: oklch(0.985 0 0)                /* White */
--primary: oklch(0.922 0 0)                   /* Very light (buttons) */
--primary-foreground: oklch(0.205 0 0)        /* Dark */
--secondary: oklch(0.269 0 0)                 /* Dark gray */
--secondary-foreground: oklch(0.985 0 0)      /* White */
--muted: oklch(0.269 0 0)                     /* Dark gray (messages) */
--muted-foreground: oklch(0.708 0 0)          /* Medium gray */
```

#### Sidebar-Specific Colors
```css
--sidebar: oklch(0.985 0 0)                   /* White in light mode */
--sidebar-foreground: oklch(0.145 0 0)        /* Dark text */
--sidebar-primary: oklch(0.205 0 0)           /* Active state */
--sidebar-accent: oklch(0.97 0 0)             /* Hover state */
```

#### Chart Colors (for future analytics)
```css
--chart-1 to --chart-5                        /* Predefined palette for 5 line/bar charts */
```

### Component Styling Approach

**ChatWindow Message Styling:**
```typescript
// User messages: right-aligned, primary background
className: `bg-primary text-primary-foreground`

// Assistant messages: left-aligned, muted background
className: `bg-muted text-muted-foreground`
```

**StepsDisplay Styling:**
```typescript
// Status icons with color coding:
// - Completed: text-green-600
// - Running: text-blue-600 (animated spin)
// - Pending: text-gray-400
// - Error: text-red-600

// Cards with left border: border-l-4 border-l-muted
// Max height for details: max-h-32 overflow-y-auto
```

**Sidebar Styling:**
```typescript
// Collapsible icon mode on mobile
// Header: border-b with logo
// Items: hover background, active state
// Footer: border-t with minimize button
// Trigger button: black background when collapsed
```

---

## 6. File Locations Summary

### Core Pages
| Path | Purpose |
|------|---------|
| `src/app/layout.tsx` | Root layout (fonts, structure) |
| `src/app/page.tsx` | Root page (redirects) |
| `src/app/globals.css` | Global styles + color system |
| `src/app/system-analysis/page.tsx` | **Main page with ChatWindow + StepsDisplay** |
| `src/app/documentation/page.tsx` | Documentation placeholder |

### Components
| Path | Purpose |
|------|---------|
| `src/components/ChatWindow.tsx` | **Chat interface for analysis** |
| `src/components/StepsDisplay.tsx` | **LLM processing steps visualization** |
| `src/components/sidebar-layout.tsx` | Layout wrapper (sidebar + content) |
| `src/components/app-sidebar.tsx` | Sidebar navigation |
| `src/components/ui/*.tsx` | 50+ shadcn components (buttons, cards, etc.) |

### Configuration
| Path | Purpose |
|------|---------|
| `tailwind.config.ts` | Tailwind configuration (minimal) |
| `next.config.ts` | Next.js configuration |
| `tsconfig.json` | TypeScript configuration |
| `biome.json` | Biome linter/formatter config |
| `src/lib/sidebar.ts` | Navigation items definition |
| `src/lib/utils.ts` | Utility functions (cn, etc.) |

---

## 7. Data Flow

### Message Handling (System Analysis)
```
User types in ChatWindow
         ↓
handleSendMessage()
         ↓
Local message added to state
         ↓
sendAnalysisQuery() API call
         ↓
Response with content + steps
         ↓
onMessagesUpdate() callback → parent state
onStepsUpdate() callback → parent state
         ↓
StepsDisplay updates (if steps exist)
ChatWindow re-renders with assistant message
```

### State Management
- **Client-side only** - No global state (Redux, Zustand, etc.)
- Each page manages its own state via `useState`
- Parent page (`system-analysis/page.tsx`) coordinates between ChatWindow and StepsDisplay
- No persistence between page navigations

---

## 8. Responsive Behavior

### Breakpoints
- **Mobile:** Sidebar converts to sheet modal (hidden by default)
- **Tablet/Desktop (`md` and up):** Sidebar always visible (collapsible to icon)

### Key Responsive Classes
```typescript
// Sidebar
group-data-[collapsible=icon]:hidden     // Hide text when collapsed
group-data-[collapsible=icon]:size-8!    // Icon-only mode

// Sidebar width variables
--sidebar-width: 256px                   // Desktop expanded
--sidebar-width-icon: 52px               // Desktop collapsed
--sidebar-width-mobile: 256px            // Mobile sheet
```

---

## 9. Potential Improvements

### 1. **Chat Window Max Width** (Large Screens)
Current: Takes full available width on desktop
```typescript
// Limit chat to readable width on ultra-wide screens
<div className="w-full max-w-4xl mx-auto">
  <ChatWindow />
</div>
```

### 2. **Resizable Panes** (Steps + Chat)
Current: Fixed 50/50 split
```typescript
// Consider adding resizable divider (use @radix-ui/primitive)
<Rnd minWidth={200} maxWidth={800}>
  <ChatWindow />
</Rnd>
```

### 3. **Sticky Header with Controls**
```typescript
// Add sticky title bar with actions
<div className="sticky top-0 bg-card border-b p-4">
  <div className="flex justify-between items-center">
    <h1>System Analysis</h1>
    <button>Export</button>
  </div>
</div>
```

### 4. **Message Time Grouping**
Current: Every message shows timestamp
```typescript
// Group messages by time/day
<div className="text-center text-xs text-muted-foreground py-2">
  Today
</div>
```

### 5. **Markdown Support in Messages**
Current: Plain text only
```typescript
// Use react-markdown for formatted responses
import ReactMarkdown from "react-markdown"
<ReactMarkdown>{msg.content}</ReactMarkdown>
```

---

## 10. Architecture Summary

```
Next.js App Router (React 18+)
├── Server Components (pages)
├── Client Components (interactive)
│   ├── ChatWindow (chat UI)
│   ├── StepsDisplay (processing steps)
│   ├── SidebarLayout (layout provider)
│   └── AppSidebar (navigation)
├── UI Library (shadcn/Radix)
├── Tailwind CSS (styling)
└── Custom Services (analysis.service.ts)
```

**Key Architectural Patterns:**
1. **Component Composition** - Small, focused components
2. **Prop Callbacks** - Parent-child communication via props
3. **Conditional Rendering** - Steps only show when needed
4. **Flexbox Layout** - Responsive without media queries
5. **CSS Variables** - Dynamic theming (light/dark mode)

---

## Conclusion

The frontend has a **well-organized, clean structure** with:
- ✅ Minimal component overlap (no z-index issues identified)
- ✅ Proper flexbox constraints for scroll areas
- ✅ Modern color system with OKLCH
- ✅ Modular component architecture
- ✅ Clear data flow and state management

**Quick Wins for Improvement:**
1. Add resizable pane divider between ChatWindow and StepsDisplay
2. Implement sticky header with export/clear buttons
3. Add markdown support for LLM responses
4. Limit chat width on ultra-wide screens
5. Add message time grouping for better UX
