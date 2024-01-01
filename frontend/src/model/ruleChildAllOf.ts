import type { RuleChildAllOfType } from "./ruleChildAllOfType";
/**
 * Generated by orval v6.23.0 🍺
 * Do not edit manually.
 * regoviz
 * api for regoviz
 * OpenAPI spec version: 1.0.0
 */
import type { RuleParent } from "./ruleParent";
import type { RuleStatement } from "./ruleStatement";

export type RuleChildAllOf = {
	parent?: RuleParent;
	statements?: RuleStatement[];
	type?: RuleChildAllOfType;
	value?: string;
};