-- name: GetIdentity :one
SELECT
    *
FROM
    device_identity
LIMIT
    1;

-- name: CreateIdentity :exec
INSERT INTO
    device_identity (
        id,
        name,
        device_type,
        os,
        os_version,
        arch,
        hostname,
        public_key,
        private_key,
        created_at
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);