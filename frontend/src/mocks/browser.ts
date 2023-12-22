import { setupWorker } from "msw/browser";
import { getDefaultMock } from "../default/default.msw";

export const worker = setupWorker(...getDefaultMock());
