import pg from "pg";

// ============================================================================
// DB CLIENT - Direct postgres access for tests that need to insert test data
// ============================================================================

let _client: pg.Client | null = null;

export async function getDB(): Promise<pg.Client> {
    if (_client) return _client;

    _client = new pg.Client({
        connectionString: process.env.DATABASE_URL,
    });

    await _client.connect();
    return _client;
}

export async function closeDB(): Promise<void> {
    if (_client) {
        await _client.end();
        _client = null;
    }
}

export async function dbQuery<T extends pg.QueryResultRow = any>(
    sql: string,
    params: unknown[] = []
): Promise<pg.QueryResult<T>> {
    const db = await getDB();
    return db.query<T>(sql, params);
}