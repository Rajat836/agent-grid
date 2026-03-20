"use client";

import { Config } from "@/base/config";
import type { RequestCreateRule } from "@/types/rules/req_create_rule";
import type { ResponseCreateRule } from "@/types/rules/res_create_rule";
import type { ResponseGetRule } from "@/types/rules/res_get_rule";
import type { ResponseRunRule } from "@/types/rules/res_run_rule";

/**
 * Rules Service
 * Handles API calls related to rules data
 */
const baseUrl = Config.apiHost;

/**
 * Creates a new rule
 * Endpoint: POST /ontology/v1/rules/create
 * @param type - The dictionary type (e.g., 'event_type', 'entity', etc.)
 * @returns Promise<ResponseCreateRule>
 */
export async function CreateRule(
  params: RequestCreateRule,
): Promise<ResponseCreateRule> {
  const url = `${baseUrl}/ontology/v1/rules/create`;

  try {
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "x-scope": "rules:create",
      },
      body: JSON.stringify(params),
    });
    if (!response.ok) {
      throw new Error(`Failed to create rule: ${response.statusText}`);
    }

    const data: ResponseCreateRule = await response.json();
    return data;
  } catch (error) {
    console.error(`Error creating rule:`, error);
    throw error;
  }
}

export async function RunRule(id: string): Promise<ResponseRunRule> {
  const url = `${baseUrl}/ontology/v1/rules/run`;

  try {
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "x-scope": "rules:manage",
      },
      body: JSON.stringify({ id }),
    });
    if (!response.ok) {
      throw new Error(`Failed to run rule: ${response.statusText}`);
    }

    const data: ResponseRunRule = await response.json();
    return data;
  } catch (error) {
    console.error(`Error running rule:`, error);
    throw error;
  }
}

export async function getRule(id: string): Promise<ResponseGetRule> {
  const url = `${baseUrl}/ontology/v1/rules/${id}`;
  console.log(`🌐 RulesService: Fetching rule with URL: ${url}`);

  try {
    const response = await fetch(url, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "x-scope": "rules:view",
      },
    });
    if (!response.ok) {
      throw new Error(
        `Failed to fetch rule with ID ${id}: ${response.statusText}`,
      );
    }
    const data: ResponseGetRule = await response.json();
    return data;
  } catch (error) {
    console.error(`Error fetching rule:`, error);
    throw error;
  }
}

export const RulesService = {
  CreateRule,
  RunRule,
  getRule,
};
