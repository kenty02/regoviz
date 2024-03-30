import { expect, test } from "@playwright/test";

test.beforeEach(async ({ page }) => {
	// Go to the starting url before each test.
	await page.goto("http://localhost:5173/");
});

test("has title", async ({ page }) => {
	// Expect a title "to contain" a substring.
	await expect(page).toHaveTitle(/regoviz/);
});
