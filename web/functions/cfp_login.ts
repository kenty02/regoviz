import { CFP_COOKIE_MAX_AGE } from "./constants";
import { getCookieKeyValue, sha256 } from "./utils";

export async function onRequestPost(context: {
	request: Request;
	env: { VITE_API_TOKEN?: string };
}): Promise<Response> {
	const { request, env } = context;
	const body = await request.formData();
	// @ts-ignore
	const { password, redirect } = Object.fromEntries(body);
	const hashedPassword = await sha256(password.toString());
	const hashedCfpPassword = await sha256(env.VITE_API_TOKEN);
	const redirectPath = redirect.toString() || "/";

	if (hashedPassword === hashedCfpPassword) {
		// Valid password. Redirect to home page and set cookie with auth hash.
		const cookieKeyValue = await getCookieKeyValue(env.VITE_API_TOKEN);

		return new Response("", {
			status: 302,
			headers: {
				"Set-Cookie": `${cookieKeyValue}; Max-Age=${CFP_COOKIE_MAX_AGE}; Path=/; HttpOnly; Secure`,
				"Cache-Control": "no-cache",
				Location: redirectPath,
			},
		});
	}
	// Invalid password. Redirect to login page with error.
	return new Response("", {
		status: 302,
		headers: {
			"Cache-Control": "no-cache",
			Location: `${redirectPath}?error=1`,
		},
	});
}
