"use client";

import { Config } from "@/base/config";

// Enable mock data for testing - set via NEXT_PUBLIC_USE_MOCK_DATA env variable
const USE_MOCK_DATA = typeof window !== "undefined" 
  ? process.env.NEXT_PUBLIC_USE_MOCK_DATA === "true"
  : false;

export interface ChatMessage {
  id: string;
  role: "user" | "assistant";
  content: string;
  timestamp: number;
  contentType?: "text" | "markdown"; // "text" for plain text, "markdown" for markdown content
}

export interface LLMStep {
  id: string;
  title: string;
  description: string;
  status: "pending" | "running" | "completed" | "error";
  startTime?: number;
  endTime?: number;
  details?: string;
}

export interface AnalysisResponse {
  message_id: string;
  content: string;
  contentType?: "text" | "markdown"; // "text" for plain text, "markdown" for markdown content
  steps: LLMStep[];
}

const baseUrl = Config.apiHost;

/**
 * Send a query to the LLM for system analysis
 * Endpoint: POST /ontology_bot/v1/analysis/query
 * 
 * If USE_MOCK_DATA is true, returns mock data for testing
 */
export async function sendAnalysisQuery(
  query: string,
): Promise<AnalysisResponse> {
  if (USE_MOCK_DATA) {
    console.log("[MOCK] Sending analysis query:", query);
    const { sendAnalysisQueryMock } = await import("./analysis.mock");
    return sendAnalysisQueryMock(query);
  }

  const url = `${baseUrl}/ontology_bot/v1/analysis/query`;

  try {
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "x-scope": "analysis:query",
      },
      body: JSON.stringify({ query }),
    });

    if (!response.ok) {
      throw new Error(`Failed to process query: ${response.statusText}`);
    }

    const data: AnalysisResponse = await response.json();
    return data;
  } catch (error) {
    console.error("Error sending analysis query:", error);
    throw error;
  }
}

/**
 * Stream LLM processing steps in real-time
 * 
 * If USE_MOCK_DATA is true, returns mock streaming data for testing
 */
export async function* streamAnalysisQuery(query: string) {
  if (USE_MOCK_DATA) {
    console.log("[MOCK] Streaming analysis query:", query);
    const { streamAnalysisQueryMock } = await import("./analysis.mock");
    yield* streamAnalysisQueryMock(query);
    return;
  }

  const url = `${baseUrl}/ontology_bot/v1/analysis/query/stream`;

  try {
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "x-scope": "analysis:query",
      },
      body: JSON.stringify({ query }),
    });

    if (!response.ok) {
      throw new Error(`Failed to process query: ${response.statusText}`);
    }

    if (!response.body) {
      throw new Error("Response body is empty");
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();

    try {
      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value);
        const lines = chunk.split("\n");

        for (const line of lines) {
          if (line.trim()) {
            try {
              const data = JSON.parse(line);
              yield data;
            } catch (e) {
              console.error("Failed to parse streaming response:", e);
            }
          }
        }
      }
    } finally {
      reader.releaseLock();
    }
  } catch (error) {
    console.error("Error streaming analysis query:", error);
    throw error;
  }
}
