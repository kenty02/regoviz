/**
 * Generated by orval v6.23.0 🍺
 * Do not edit manually.
 * regoviz
 * api for regoviz
 * OpenAPI spec version: 1.0.0
 */

export type GetCallTreeParams = {
	/**
	 * The sample name to analyze
	 */
	sampleName: string;
	/**
	 * The entrypoint rule to analyze
	 */
	entrypoint: string;
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
	query?: string;
};
