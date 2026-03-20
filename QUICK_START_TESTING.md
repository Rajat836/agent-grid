# Quick Start: Test Chat with Agent Steps 🚀

## Step 1: Start the Frontend (90 seconds)

```bash
cd /home/abhi/workspace/fyscal/ontology_bot/panel
npm run dev
```

You should see:
```
➜  Local:        http://localhost:3000/system-analysis
```

## Step 2: Open in Browser

Navigate to: **http://localhost:3000/system-analysis**

You should see:
- Left sidebar with navigation
- "System Analysis" header
- Chat window labeled "LIVE CONVERSATION"
- Empty message area with chat icon

## Step 3: Test the Chat

### Type a message and hit Enter

Try any of these:
- "Analyze my system performance"
- "What are the current bottlenecks?"
- "How can I optimize my server?"
- "Show me the system metrics"

## What You'll See

### Immediately:
✅ Your message appears in a blue bubble on the right
✅ A "Processing..." indicator shows while the agent works

### After ~2 seconds:
✅ The "Processing Steps" panel appears below the chat
✅ Step 1: "Parsing Query" - ✓ Complete (green checkmark)
✅ Step 2: "Gathering System Information" - ✓ Complete
✅ Step 3: "Analyzing Patterns" - ✓ Complete  
✅ Step 4: "Generating Recommendations" - ✓ Complete

### Final Result:
✅ Agent response appears in a teal bubble on the left
✅ Each step shows:
  - Status icon (green checkmark for completed)
  - Title and description
  - Timing information (e.g., "1.2s")
  - Details panel with specific information

## Visual Verification Checklist

As you test, verify these visual elements:

**Chat Window:**
- [ ] User message is blue with white text
- [ ] User message appears on the right
- [ ] Assistant message is teal/emerald with darker text
- [ ] Assistant message appears on the left
- [ ] Messages show timestamps
- [ ] Chat input box at bottom with send button
- [ ] Send button is blue/indigo gradient

**Steps Display:**
- [ ] "Processing Steps" header appears
- [ ] Steps appear in order (1 → 4)
- [ ] Completed steps have ✓ green icon
- [ ] Each step shows title and description
- [ ] Timing shows (e.g., "0.5s")
- [ ] Details section expands with information
- [ ] Left border is colored (purple for section)

**Overall:**
- [ ] No blur or opacity issues
- [ ] Text is crisp and readable
- [ ] Colors are vibrant and distinct
- [ ] Layout doesn't overlap or shift
- [ ] Sidebar stays visible on left

## Features to Test

### 1. Multiple Messages
Send several messages in sequence:
```
1. "Analyze my system"
2. "Why is CPU high?"
3. "Give me optimization tips"
```

Verify:
- All messages stay in history
- Each query gets fresh steps
- Chat scrolls automatically

### 2. Keyboard Input
- Type and press **Enter** to send
- Try clicking the send button

### 3. Real-time Step Display
Watch the steps appear in real-time:
- Parse → Complete (1s)
- Gather → Running (1s) 
- Analyze → Running (2s)
- Generate → Complete (1s)
- Message → Appears

### 4. Error Handling (Optional)
Try intentionally triggering errors:
- Send empty message (should disable send button)
- Rapid multiple submissions (should queue)

## Connection Verification

### Check Browser Console (F12 → Console)

Look for these logs:
```
[MOCK] Sending analysis query: "Your question here"
```

This confirms the mock data is being used.

### Check Network Tab (F12 → Network)

**Expected:** NO network requests to `/ontology_bot/v1/analysis/query`

This proves mock data is working locally.

## Current Configuration

Your `.env.local` is set to:
```
NEXT_PUBLIC_USE_MOCK_DATA=true
```

This means all chat queries use realistic mock data, perfect for testing the UI and workflow.

## Next Steps After Testing Mock

Once you've verified everything works:

### Option 1: Test with Real Backend
1. Implement backend endpoint (see `TESTING_CHAT_STEPS.md` Option 2)
2. Update `.env.local`:
   ```
   NEXT_PUBLIC_USE_MOCK_DATA=false
   ```
3. Restart frontend
4. Same testing process, now with real data

### Option 2: Keep Using Mock for Development
Keep current setup - mock provides realistic data without backend dependency.

## Files Added/Modified

**New Files:**
- `panel/src/services/analysis.mock.ts` - Mock data service
- `TESTING_CHAT_STEPS.md` - Detailed testing guide
- `QUICK_START_TESTING.md` - This file

**Modified Files:**
- `panel/src/services/analysis.service.ts` - Added mock support
- `panel/.env.local` - Added NEXT_PUBLIC_USE_MOCK_DATA=true

## Estimated Testing Time

- Initial load: 2-3 minutes
- First query test: 30 seconds
- Full feature test: 5 minutes
- Total: ~10 minutes to fully verify

## Success Indicators

You'll know it's working when:
1. ✅ Chat message appears immediately
2. ✅ "Processing Steps" panel shows up
3. ✅ All 4 steps display with completion
4. ✅ Agent response appears
5. ✅ Multiple queries work sequentially
6. ✅ No console errors

Ready to test? Start with Step 1 above! 🎉
