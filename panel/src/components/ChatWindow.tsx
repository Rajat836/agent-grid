"use client";

import React, { useRef, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Send, Loader2 } from "lucide-react";
import type { ChatMessage } from "@/services/analysis.service";
import { sendAnalysisQuery } from "@/services/analysis.service";
import { MarkdownRenderer } from "@/components/MarkdownRenderer";

interface ChatWindowProps {
  onMessagesUpdate?: (messages: ChatMessage[]) => void;
  onStepsUpdate?: (steps: any[]) => void;
}

export function ChatWindow({ onMessagesUpdate, onStepsUpdate }: ChatWindowProps) {
  const [messages, setMessages] = React.useState<ChatMessage[]>([]);
  const [input, setInput] = React.useState("");
  const [isLoading, setIsLoading] = React.useState(false);
  const scrollRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages]);

  const handleSendMessage = async () => {
    if (!input.trim()) return;

    // Add user message
    const userMessage: ChatMessage = {
      id: `msg-${Date.now()}`,
      role: "user",
      content: input,
      timestamp: Date.now(),
    };

    const newMessages = [...messages, userMessage];
    setMessages(newMessages);
    onMessagesUpdate?.(newMessages);
    setInput("");
    setIsLoading(true);

    try {
      // Call LLM service
      const response = await sendAnalysisQuery(input);

      // Notify parent of steps
      onStepsUpdate?.(response.steps);

      // Add assistant message with markdown support
      const assistantMessage: ChatMessage = {
        id: `msg-${Date.now()}`,
        role: "assistant",
        content: response.content,
        contentType: response.contentType || "text",
        timestamp: Date.now(),
      };

      const updatedMessages = [...newMessages, assistantMessage];
      setMessages(updatedMessages);
      onMessagesUpdate?.(updatedMessages);
    } catch (error) {
      console.error("Error sending message:", error);

      // Add error message
      const errorMessage: ChatMessage = {
        id: `msg-${Date.now()}`,
        role: "assistant",
        content: "Sorry, there was an error processing your query. Please try again.",
        timestamp: Date.now(),
      };

      const updatedMessages = [...newMessages, errorMessage];
      setMessages(updatedMessages);
      onMessagesUpdate?.(updatedMessages);
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  return (
    <div className="flex flex-col gap-3 border border-gray-300 dark:border-gray-700 rounded-lg p-4 bg-white dark:bg-slate-900 shadow-sm transition-shadow h-full overflow-hidden">
      {/* Chat Header */}
      <div className="shrink-0 pb-3 border-b border-gray-300 dark:border-gray-700 flex items-center gap-2">
        <div className="w-2 h-2 bg-gray-400 rounded-full animate-pulse"></div>
        <span className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider">Live Conversation</span>
      </div>

      {/* Chat Messages */}
      <ScrollArea className="flex-1 min-h-0 border border-gray-300 dark:border-gray-700 rounded-lg p-4 bg-gray-100 dark:bg-slate-800">
        <div ref={scrollRef} className="space-y-3">
          {messages.length === 0 ? (
            <div className="text-center text-muted-foreground text-sm py-12 space-y-3">
              <div className="text-4xl">💬</div>
              <p className="font-medium">Start a conversation</p>
              <p className="text-xs">Ask questions about system analysis to get started</p>
            </div>
          ) : (
            messages.map((msg) => (
              <div key={msg.id} className={`flex animate-slide-up ${msg.role === "user" ? "justify-end" : "justify-start"}`}>
                <div
                  className={`max-w-2xl px-4 py-3 rounded-lg transition-all ${
                    msg.role === "user"
                      ? "bg-accent-blue dark:bg-blue-600 text-white shadow-sm font-medium text-sm"
                      : "bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 border border-gray-300 dark:border-gray-700"
                  }`}
                >
                  {msg.role === "user" ? (
                    <p>{msg.content}</p>
                  ) : msg.contentType === "markdown" ? (
                    <MarkdownRenderer content={msg.content} className="text-gray-900 dark:text-gray-100" />
                  ) : (
                    <p className="text-sm leading-relaxed whitespace-pre-wrap">{msg.content}</p>
                  )}
                  <time className="text-xs opacity-60 mt-2 block">
                    {new Date(msg.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                  </time>
                </div>
              </div>
            ))
          )}
          {isLoading && (
            <div className="flex justify-start animate-slide-up">
              <div className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 px-4 py-2 rounded-lg flex items-center gap-2 border border-gray-300 dark:border-gray-700">
                <Loader2 className="w-4 h-4 animate-spin text-gray-600 dark:text-gray-400" />
                <span className="text-sm font-medium">Processing your request...</span>
              </div>
            </div>
          )}
        </div>
      </ScrollArea>

      {/* Input Section */}
      <div className="flex gap-2 shrink-0 border-t border-gray-300 dark:border-gray-700 pt-3">
        <Input
          placeholder="Ask a question about system analysis..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyPress={handleKeyPress}
          disabled={isLoading}
          className="flex-1 border-gray-300 dark:border-gray-700 bg-white/50 dark:bg-slate-800/50 text-gray-900 dark:text-gray-100 focus:border-accent-blue dark:focus:border-blue-400 focus:ring-2 focus:ring-accent-blue/20 dark:focus:ring-blue-400/20 placeholder:text-gray-500 dark:placeholder:text-gray-400"
        />
        <Button
          onClick={handleSendMessage}
          disabled={isLoading || !input.trim()}
          size="sm"
          className="bg-accent-coral hover:bg-red-600 dark:bg-accent-coral dark:hover:bg-red-600 text-white font-medium transition-all"
        >
          {isLoading ? (
            <Loader2 className="w-4 h-4 animate-spin" />
          ) : (
            <Send className="w-4 h-4" />
          )}
        </Button>
      </div>
    </div>
  );
}
