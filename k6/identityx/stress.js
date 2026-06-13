import { check, sleep } from 'k6';
import http from 'k6/http';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

// =============================================================
// CONFIG
// =============================================================

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

const JSON_HEADERS = { 'Content-Type': 'application/json' };

function authHeaders(token) {
    return { headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' } };
}

// =============================================================
// SCENARIO
// =============================================================

export const options = {
    scenarios: {
        breaking_point: {
            executor: 'ramping-arrival-rate',
            startRate: 5,
            timeUnit: '1s',
            stages: [
                { duration: '1m', target: 20 },   // warm up
                { duration: '2m', target: 50 },   // light load
                // { duration: '2m', target: 150 },  // moderate
                // { duration: '2m', target: 300 },  // heavy
                // { duration: '2m', target: 500 },  // stress
                // { duration: '2m', target: 800 },  // breaking point
                { duration: '1m', target: 0 },    // cooldown — watch recovery
            ],
            preAllocatedVUs: 50,
            maxVUs: 50,
        },
    },
    // no thresholds on stress — we want to observe natural failure, not abort early
};

export function setup() {
    const res = http.post(
        `${BASE_URL}/auth/setup`,
        JSON.stringify({ email: 'admin@trieoh.com', password: 'S3cretPa$$' }),
        { headers: JSON_HEADERS }
    );
    // 200 = first time, 409/4xx = already set up, both are fine
    console.log(`setup: ${res.status}`);
}

export function teardown() {
    console.log('stress test complete');
    // optionally hit a cleanup endpoint if you add one
}

// =============================================================
// HELPERS
// =============================================================

function login(email, password) {
    const res = http.post(
        `${BASE_URL}/auth/login`,
        JSON.stringify({ email, password }),
        { headers: JSON_HEADERS }
    );
    check(res, { 'login 200': r => r.status === 200 });
    if (res.status !== 200) return null;
    return {
        access_token: res.json('data.access_token'),
        refresh_token: res.json('data.refresh_token'),
    };
}

function introspect(token) {
    const res = http.get(`${BASE_URL}/auth/introspect`, authHeaders(token));
    check(res, { 'introspect 200': r => r.status === 200 });
}

function createOrg(token, slug) {
    const res = http.post(
        `${BASE_URL}/organizations`,
        JSON.stringify({ name: `Org ${slug}`, slug }),
        authHeaders(token)
    );
    check(res, { 'create org 2xx': r => r.status >= 200 && r.status < 300 });
    if (res.status >= 200 && res.status < 300) return res.json('data.id');
    return null;
}

function listOrgs(token) {
    const res = http.get(`${BASE_URL}/organizations`, authHeaders(token));
    check(res, { 'list orgs 200': r => r.status === 200 });
}

function addMember(token, orgId, email) {
    const res = http.post(
        `${BASE_URL}/organizations/${orgId}/members`,
        JSON.stringify({ actor_email: email, role: 'member' }),
        authHeaders(token)
    );
    check(res, { 'add member 2xx': r => r.status >= 200 && r.status < 300 });
}

function listMembers(token, orgId) {
    const res = http.get(`${BASE_URL}/organizations/${orgId}/members`, authHeaders(token));
    check(res, { 'list members 200': r => r.status === 200 });
}

function refreshToken(refreshTok) {
    const res = http.post(
        `${BASE_URL}/auth/refresh`,
        null,
        { headers: { refresh_token: refreshTok } }
    );
    check(res, { 'refresh 200': r => r.status === 200 });
    if (res.status !== 200) return null;
    return res.json('data.access_token');
}

function logout(accessTok, refreshTok) {
    const res = http.post(
        `${BASE_URL}/auth/logout`,
        null,
        { headers: { access_token: accessTok, refresh_token: refreshTok } }
    );
    check(res, { 'logout 2xx': r => r.status >= 200 && r.status < 300 });
}

// =============================================================
// MAIN FLOW — each VU iteration is a full user session
// =============================================================

export default function () {
    // unique per-VU identity so users don't collide
    const uid = randomString(8).toLowerCase();
    const email = `stress_${uid}@trieoh.com`;
    const password = 'S3cretPa$$';

    // 1. register a fresh user
    const reg = http.post(
        `${BASE_URL}/auth/register`,
        JSON.stringify({ email, password }),
        { headers: JSON_HEADERS }
    );
    check(reg, { 'register 2xx': r => r.status >= 200 && r.status < 300 });
    if (reg.status >= 300) return; // bail if registration failed

    sleep(0.5);

    // 2. login
    const tokens = login(email, password);
    if (!tokens) return;

    sleep(0.3);

    // 3. introspect (common read — hits jwt validation + spicedb)
    introspect(tokens.access_token);

    sleep(0.2);

    // 4. create org
    const slug = `org-${uid}`;
    const orgId = createOrg(tokens.access_token, slug);

    sleep(0.2);

    // 5. list orgs
    listOrgs(tokens.access_token);

    sleep(0.2);

    if (orgId) {
        // 6. add a member
        const memberEmail = `member_${uid}@trieoh.com`;
        addMember(tokens.access_token, orgId, memberEmail);

        sleep(0.2);

        // 7. list members
        listMembers(tokens.access_token, orgId);

        sleep(0.2);
    }

    // 8. refresh token (exercises token rotation path)
    const newToken = refreshToken(tokens.refresh_token);

    sleep(0.3);

    // 9. logout with whichever token we have
    logout(newToken || tokens.access_token, tokens.refresh_token);
}