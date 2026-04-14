CREATE TABLE
    devices (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        public_key TEXT NOT NULL,
        state TEXT NOT NULL,
        last_seen INTEGER NOT NULL,
        trusted_at INTEGER NOT NULL
    );

CREATE TABLE
    device_identity (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        device_type TEXT NOT NULL,
        os TEXT NOT NULL,
        os_version TEXT NOT NULL,
        arch TEXT NOT NULL,
        hostname TEXT NOT NULL,
        public_key TEXT NOT NULL,
        private_key TEXT NOT NULL,
        created_at INTEGER NOT NULL
    );