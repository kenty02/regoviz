import type { RuleChildAllOfType } from "./ruleChildAllOfType";
/**
 * Generated by orval v6.23.0 🍺
 * Do not edit manually.
 * regoviz
 * api for regoviz
 * OpenAPI spec version: 1.0.0
 */
import type { RuleStatement } from "./ruleStatement";

export type RuleChildAllOf = {
	statements: RuleStatement[];
	type: RuleChildAllOfType;
	value: string;
};
