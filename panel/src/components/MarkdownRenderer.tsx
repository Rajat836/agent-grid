"use client";

import React from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

interface MarkdownRendererProps {
  content: string;
  className?: string;
}

export function MarkdownRenderer({ content, className = "" }: MarkdownRendererProps) {
  return (
    <div className={`prose prose-sm dark:prose-invert max-w-none ${className}`}>
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={{
          // Style lists
          ul: ({ node, ...props }) => (
            <ul className="list-disc list-inside space-y-1 my-2" {...props} />
          ),
          ol: ({ node, ...props }) => (
            <ol className="list-decimal list-inside space-y-1 my-2" {...props} />
          ),
          li: ({ node, ...props }) => (
            <li className="text-sm" {...props} />
          ),
          
          // Style headings
          h1: ({ node, ...props }) => (
            <h1 className="text-lg font-bold my-2" {...props} />
          ),
          h2: ({ node, ...props }) => (
            <h2 className="text-base font-bold my-2" {...props} />
          ),
          h3: ({ node, ...props }) => (
            <h3 className="text-sm font-bold my-2" {...props} />
          ),
          
          // Style code blocks
          code: (({ node, inline, ...props }: any) => {
            if (inline) {
              return (
                <code
                  className="bg-gray-200 dark:bg-gray-700 px-1.5 py-0.5 rounded font-mono text-xs"
                  {...props}
                />
              );
            }
            return (
              <code
                className="block bg-gray-900 dark:bg-black text-gray-100 p-3 rounded-lg font-mono text-xs overflow-x-auto my-2"
                {...props}
              />
            );
          }) as any,
          pre: ({ node, ...props }) => (
            <pre className="my-2 overflow-x-auto" {...props} />
          ),
          
          // Style tables
          table: ({ node, ...props }) => (
            <table className="border-collapse border border-gray-300 dark:border-gray-600 text-sm my-2" {...props} />
          ),
          thead: ({ node, ...props }) => (
            <thead className="bg-gray-100 dark:bg-gray-800" {...props} />
          ),
          th: ({ node, ...props }) => (
            <th className="border border-gray-300 dark:border-gray-600 px-2 py-1 text-left font-bold" {...props} />
          ),
          td: ({ node, ...props }) => (
            <td className="border border-gray-300 dark:border-gray-600 px-2 py-1" {...props} />
          ),
          
          // Style blockquotes
          blockquote: ({ node, ...props }) => (
            <blockquote
              className="border-l-4 border-gray-400 dark:border-gray-600 pl-3 italic text-gray-600 dark:text-gray-400 my-2"
              {...props}
            />
          ),
          
          // Style links
          a: ({ node, ...props }) => (
            <a
              className="text-blue-600 dark:text-blue-400 underline hover:opacity-80"
              target="_blank"
              rel="noopener noreferrer"
              {...props}
            />
          ),
          
          // Style paragraphs
          p: ({ node, ...props }) => (
            <p className="text-sm my-1 leading-relaxed" {...props} />
          ),

          // Style horizontal rules
          hr: ({ node, ...props }) => (
            <hr className="my-2 border-gray-300 dark:border-gray-600" {...props} />
          ),

          // Style strong
          strong: ({ node, ...props }) => (
            <strong className="font-bold" {...props} />
          ),

          // Style emphasis
          em: ({ node, ...props }) => (
            <em className="italic" {...props} />
          ),
        }}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
}
