export enum RuleType {
  SIMPLE = "simple",
  SCENARIO = "scenario",
}

export enum ConditionType {
  AND = "and",
  OR = "or",
  NOT = "not",
}

export enum SubType {
  transaction = "transaction_subtype",
  action = "action_subtype",
  interval = "interval",
}

export enum EventType {
  transaction = "transaction",
  action = "action",
}

export const EventSubTypeMapping = {
  transaction: SubType.transaction,
  action: SubType.action,
  interval: SubType.interval,
};

export const RuleTypeLabels: Record<RuleType, string> = {
  [RuleType.SIMPLE]: "Simple Rule",
  [RuleType.SCENARIO]: "Scenario Based Rule",
};

export const DicEventColumnMapping: Record<EventType, string> = {
  [EventType.transaction]: "txn_events",
  [EventType.action]: "action_events",
};

export const getInputConfig = (fieldType: string | undefined) => {
  switch (fieldType) {
    case "number":
      return { type: "number", placeholder: "Enter number" };
    case "date":
      return { type: "date", placeholder: "Select date" };
    default:
      return { type: "text", placeholder: "Enter value" };
  }
};

export enum Operators {
  isNull = "is_null",
  isNotNull = "is_not_null",
  between = "between",
  in = "in",
  notIn = "not_in",
}
