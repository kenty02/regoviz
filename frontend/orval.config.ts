import { defineConfig } from "orval";

export default defineConfig({
	regoviz: {
		input: "../openapi.yml",
		output: {
			mode: "tags-split",
			target: "./src/regoviz-client.ts",
			schemas: "./src/model",
			client: "react-query",
			// https://github.com/anymaniax/orval/issues/1119
			mock: false,
			override: {
				query: {
					useSuspenseQuery: true,
					useSuspenseInfiniteQuery: true,
				},
			},
		},
		hooks: {
			// generated file paths are prefixed with this
			afterAllFilesWrite: "npx @biomejs/biome check --apply-unsafe",
		},
	},
});
