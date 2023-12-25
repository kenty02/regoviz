/**
 * Shows how to restrict access using the HTTP Basic schema.
 * @see https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication
 * @see https://tools.ietf.org/html/rfc7617
 *
 * A user-id containing a colon (":") character is invalid, as the
 * first colon in a user-pass string separates user and password.
 */
const BASIC_USER = "user";

async function errorHandling(context) {
	try {
		return await context.next();
	} catch (err) {
		return new Response(`${err.message}\n${err.stack}`, { status: 500 });
	}
}

async function handleRequest({ next, request, env }) {
	// todo: 認証が簡単すぎるので直したい
	const validPass = env.VITE_API_TOKEN;
	if (validPass == null || validPass === "") {
		return new Response("You need to set API_TOKEN.", {
			status: 500,
		});
	}

	// The "Authorization" header is sent when authenticated.
	if (request.headers.has("Authorization")) {
		const Authorization = request.headers.get("Authorization");
		// Throws exception when authorization fails.

		const [scheme, encoded] = Authorization.split(" ");

		// The Authorization header must start with Basic, followed by a space.
		if (!encoded || scheme !== "Basic") {
			return new Response("The Authorization header must start with Basic", {
				status: 400,
			});
		}

		// Decodes the base64 value and performs unicode normalization.
		// @see https://datatracker.ietf.org/doc/html/rfc7613#section-3.3.2 (and #section-4.2.2)
		// @see https://dev.mozilla.org/docs/Web/JavaScript/Reference/Global_Objects/String/normalize
		const buffer = Uint8Array.from(atob(encoded), (character) =>
			character.charCodeAt(0),
		);
		const decoded = new TextDecoder().decode(buffer).normalize();

		// The username & password are split by the first colon.
		//=> example: "username:password"
		const index = decoded.indexOf(":");

		// The user & password are split by the first colon and MUST NOT contain control characters.
		// @see https://tools.ietf.org/html/rfc5234#appendix-B.1 (=> "CTL = %x00-1F / %x7F")
		if (index === -1 || /[\0-\x1F\x7F]/.test(decoded)) {
			return new Response("Invalid authorization value.", { status: 400 });
		}

		const user = decoded.substring(0, index);
		const pass = decoded.substring(index + 1);

		if (BASIC_USER !== user) {
			return new Response("Invalid credentials.", { status: 401 });
		}

		if (validPass !== pass) {
			return new Response("Invalid credentials.", { status: 401 });
		}

		// Only returns this response when no exception is thrown.
		return await next();
	}

	// Not authenticated.
	return new Response("You need to login.", {
		status: 401,
		headers: {
			// Prompts the user for credentials.
			"WWW-Authenticate": 'Basic realm="my scope", charset="UTF-8"',
		},
	});
}

export const onRequest = [errorHandling, handleRequest];
