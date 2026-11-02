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
	// Profile handles are free text and may contain a dot.
	if (pathname.startsWith("/profiles/by-handle/")) return false;
	// A file extension on the last segment: index.php, gcp-key.json, ...
	const last = pathname.slice(pathname.lastIndexOf("/") + 1);
	return /\.[a-z0-9]{1,8}$/i.test(last);
}

export function notFound(): Response {
	return new Response("Not Found", { status: 404 });
}
