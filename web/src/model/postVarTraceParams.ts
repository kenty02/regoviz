/**
 * Generated by orval v6.23.0 🍺
 * Do not edit manually.
 * regoviz
 * api for regoviz
 * OpenAPI spec version: 1.0.0
 */

export type PostVarTraceParams = {
	/**
	 * The rego code to analyze
	 */
	policy: string;
	/**
	 * The commands to analyze
	 */
	commands: string;
	/**
	 * The input to policy
	 */
	input?: string;
	/**
	 * The data to policy
	 */
	data?: string;
	/**
	 * The query to policy
	 */
	query: string;
};
