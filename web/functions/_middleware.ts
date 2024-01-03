import { CFP_ALLOWED_PATHS } from "./constants";
import { getTemplate } from "./template";
import { getCookieKeyValue } from "./utils";

export async function onRequest(context: {
	request: Request;
	next: () => Promise<Response>;
	env: { VITE_API_TOKEN?: string };
}): Promise<Response> {
	const { request, next, env } = context;
	const { pathname, searchParams } = new URL(request.url);
	const { error } = Object.fromEntries(searchParams);
	const cookie = request.headers.get("cookie") || "";
	const cookieKeyValue = await getCookieKeyValue(env.VITE_API_TOKEN);

	if (
		cookie.includes(cookieKeyValue) ||
		CFP_ALLOWED_PATHS.includes(pathname) ||
		!env.VITE_API_TOKEN
	) {
		// Correct hash in cookie, allowed path, or no password set.
		// Continue to next middleware.
		return await next();
	}
	// No cookie or incorrect hash in cookie. Redirect to login.
	return new Response(
		getTemplate({ redirectPath: pathname, withError: error === "1" }),
		{
			headers: {
				"content-type": "text/html",
			},
		},
	);
}
