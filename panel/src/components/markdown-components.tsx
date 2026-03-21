"use client";

import type { ReactNode } from "react";
import { MermaidDiagram } from "@/components/MermaidDiagram";

function normalizeCodeContent(children: ReactNode) {
  if (Array.isArray(children)) {
    return children.join("");
  }

  return typeof children === "string" ? children : String(children ?? "");
}

export const markdownComponents = {
  code: ({ className, children, ...props }: any) => {
    const language = className?.replace("language-", "").trim().toLowerCase();
    const codeContent = normalizeCodeContent(children).replace(/\n$/, "");

    if (language === "mermaid") {
      return <MermaidDiagram chart={codeContent} />;
    }

    return className ? (
      <pre>
        <code className={className} {...props}>
          {children}
        </code>
      </pre>
    ) : (
      <code {...props}>{children}</code>
    );
  },
};
