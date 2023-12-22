import {defineConfig} from "orval";

export default defineConfig({
	regoviz: {
		input: "../openapi.yml",
		output: {
			mode: "tags-split",
			target: "./src/regoviz-client.ts",
			schemas: "./src/model",
			client: "react-query",
			mock: true,
			override: {
				query: {
					useSuspenseQuery: true,
					useSuspenseInfiniteQuery: true,
				}
			}
		},
		hooks: {
			afterAllFilesWrite: "npx @biomejs/biome format --write ./src",
		},
	},
});
