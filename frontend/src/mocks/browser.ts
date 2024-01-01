import { setupWorker } from "msw/browser";
// import { getDefaultMock } from "../default/default.msw";
const getDefaultMock = (): [] => {
	throw new Error(
		"Mock disabled until https://github.com/anymaniax/orval/issues/1119 fixes",
	);
};

export const worker = setupWorker(...getDefaultMock());
