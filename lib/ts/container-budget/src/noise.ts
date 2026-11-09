// Scanner traffic (/.env, /index.php, /wp-login, /) never matches an API
// route, but reaching the container would still wake it. No API path has a
// file extension or a dot segment except these.
const ALLOWED = new Set(["/.well-known/jwks.json", "/docs/openapi.yml"]);

// True when pathname cannot be an API route, so the Worker can answer it
// without starting a container.
export function isNoise(pathname: string): boolean {
	if (ALLOWED.has(pathname)) return false;
	if (pathname === "/" || pathname === "") return true;
	// Dotfiles and dot directories: /.env, /.git/config, /api/.aws/...
	if (pathname.includes("/.")) return true;
	// Extensionless scanner targets: Spring actuator, CGI scripts, PHP info and
	// env dumps, GitLab sign-in and version probes.
	if (/^\/(actuator|cgi-bin)(\/|$)/.test(pathname)) return true;
	if (/^\/(phpinfo|api\/env|users\/sign_in|api\/v4\/version)$/.test(pathname)) return true;
	// Profile handles are free text and may contain a dot.
	if (pathname.startsWith("/profiles/by-handle/")) return false;
	// A file extension on the last segment: index.php, gcp-key.json, ...
	const last = pathname.slice(pathname.lastIndexOf("/") + 1);
	return /\.[a-z0-9]{1,8}$/i.test(last);
}

// True when pathname can reach one of the API's handlers: it is not noise,
// and its first segment is one the API serves (`routes`, generated from its
// spec into cloudflare/src/routes.ts). Scanners probing paths no API has,
// like /user, /graphql or /%2eenv, get 404 without waking the container.
export function isRoute(pathname: string, routes: ReadonlySet<string>): boolean {
	if (ALLOWED.has(pathname)) return true;
	if (isNoise(pathname)) return false;
	const end = pathname.indexOf("/", 1);
	return routes.has(pathname.slice(1, end === -1 ? undefined : end));
}

export function notFound(): Response {
	return new Response("Not Found", { status: 404 });
}
