"use client";

import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CheckCircle2, Clock, AlertCircle, Loader2 } from "lucide-react";
import type { LLMStep } from "@/services/analysis.service";

interface StepsDisplayProps {
  steps: LLMStep[];
  isLoading?: boolean;
}

export function StepsDisplay({ steps, isLoading }: StepsDisplayProps) {
  if (steps.length === 0 && !isLoading) {
    return null;
  }

  return (
    <div className="space-y-4">
      <div className="space-y-3">
        {steps.map((step, index) => (
          <Card key={step.id} className={`border-l-4 transition-all hover:shadow-lg bg-white dark:bg-slate-800 border border-gray-300 dark:border-gray-700 ${
            step.status === "completed" 
              ? "border-l-green-600 dark:border-l-green-500" 
              : step.status === "running" 
              ? "border-l-blue-600 dark:border-l-blue-500" 
              : step.status === "error" 
              ? "border-l-red-600 dark:border-l-red-500" 
              : "border-l-amber-600 dark:border-l-amber-500"
          }`}>
            <CardContent className="p-4">
              <div className="flex items-start gap-3">
                {/* Status Icon with enhanced styling */}
                <div className="mt-0.5 flex-shrink-0">
                  {step.status === "completed" && (
                    <div className="p-1 bg-green-100 dark:bg-green-900/30 rounded-full">
                      <CheckCircle2 className="w-5 h-5 text-green-600 dark:text-green-400" />
                    </div>
                  )}
                  {step.status === "running" && (
                    <div className="p-1 bg-blue-100 dark:bg-blue-900/30 rounded-full">
                      <Loader2 className="w-5 h-5 text-blue-600 dark:text-blue-400 animate-spin" />
                    </div>
                  )}
                  {step.status === "pending" && (
                    <div className="p-1 bg-amber-100 dark:bg-amber-900/30 rounded-full">
                      <Clock className="w-5 h-5 text-amber-600 dark:text-amber-400" />
                    </div>
                  )}
                  {step.status === "error" && (
                    <div className="p-1 bg-red-100 dark:bg-red-900/30 rounded-full">
                      <AlertCircle className="w-5 h-5 text-red-600 dark:text-red-400" />
                    </div>
                  )}
                </div>

                {/* Step Details */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between gap-2">
                    <h4 className="font-semibold text-sm text-gray-900 dark:text-gray-100">{step.title}</h4>
                    {step.startTime && step.endTime && (
                      <span className="text-xs font-medium text-gray-600 dark:text-gray-400">
                        {((step.endTime - step.startTime) / 1000).toFixed(2)}s
                      </span>
                    )}
                  </div>
                  <p className="text-xs text-gray-600 dark:text-gray-400 mt-1.5 leading-relaxed">
                    {step.description}
                  </p>
                  {step.details && (
                    <div className="mt-3 p-3 bg-gray-100 dark:bg-slate-700 rounded-lg text-xs text-gray-900 dark:text-gray-100 whitespace-pre-wrap break-words max-h-40 overflow-y-auto border border-gray-300 dark:border-gray-600 font-mono">
                      {step.details}
                    </div>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>
        ))}

        {isLoading && (
          <Card className="border-l-4 border-l-blue-600 dark:border-l-blue-500 bg-white dark:bg-slate-800 border border-gray-300 dark:border-gray-700">
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="p-1 bg-blue-100 dark:bg-blue-900/30 rounded-full flex-shrink-0">
                  <Loader2 className="w-5 h-5 text-blue-600 dark:text-blue-400 animate-spin" />
                </div>
                <div>
                  <h4 className="font-semibold text-sm text-gray-900 dark:text-gray-100">Processing...</h4>
                  <p className="text-xs text-gray-600 dark:text-gray-400">
                    LLM is analyzing your request, please wait
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
