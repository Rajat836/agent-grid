"use client";

import React from "react";
import { SidebarLayout } from "@/components/sidebar-layout";
import { ChatWindow } from "@/components/ChatWindow";
import { StepsDisplay } from "@/components/StepsDisplay";
import type { ChatMessage, LLMStep } from "@/services/analysis.service";

export default function SystemAnalysis() {
  const [steps, setSteps] = React.useState<LLMStep[]>([]);
  const [messages, setMessages] = React.useState<ChatMessage[]>([]);
  const [isLoading, setIsLoading] = React.useState(false);

  const handleMessagesUpdate = (newMessages: ChatMessage[]) => {
    setMessages(newMessages);
    // Check if the last message is from assistant
    if (newMessages.length > 0 && newMessages[newMessages.length - 1].role === "assistant") {
      setIsLoading(false);
    }
  };

  const handleStepsUpdate = (newSteps: LLMStep[]) => {
    setSteps(newSteps);
    setIsLoading(true);
  };

  return (
    <SidebarLayout>
      <div className="flex flex-col h-full gap-4 px-6 py-4 bg-gray-100 dark:bg-slate-950">
        {/* Header */}
        <div className="shrink-0 border-b border-gray-300 pb-4 rounded-lg p-4 bg-white dark:bg-slate-900">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
            System Analysis
          </h1>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">Ask questions and analyze system behavior</p>
        </div>

        <div className="flex flex-col min-h-0 gap-4 flex-1 overflow-hidden">
          {/* Chat Window - Fixed at top with 60% or adaptive height */}
          <div className="shrink-0 h-96 md:h-[500px] min-h-80">
            <ChatWindow onMessagesUpdate={handleMessagesUpdate} onStepsUpdate={handleStepsUpdate} />
          </div>

          {/* Steps Display - Below chat, scrollable, only shows when needed */}
          {(steps.length > 0 || isLoading) && (
            <div className="flex-1 min-h-0 border border-gray-300 rounded-lg p-4 bg-white dark:bg-slate-900 overflow-y-auto shadow-sm">
              <div className="space-y-2 mb-2">
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 bg-gray-400 rounded-full animate-pulse"></div>
                  <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                    Processing Steps
                  </h2>
                </div>
                <div className="h-1 w-12 bg-gray-300 rounded-full"></div>
              </div>
              <StepsDisplay steps={steps} isLoading={isLoading} />
            </div>
          )}
        </div>
      </div>
    </SidebarLayout>
  );
}
