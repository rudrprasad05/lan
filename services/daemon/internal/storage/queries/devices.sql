-- name: GetDeviceByID :one
SELECT
    *
FROM
    devices
WHERE
    id = ?;

-- name: UpsertDevice :exec
INSERT INTO
    devices (
        id,
        name,
        public_key,
        state,
        last_seen,
        trusted_at
    )
VALUES
    (?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
    name = excluded.name,
    public_key = excluded.public_key,
    state = excluded.state,
    last_seen = excluded.last_seen,
    trusted_at = excluded.trusted_at;

-- name: UpdateTrustedAt :exec
UPDATE devices
SET
    trusted_at = ?
WHERE
    id = ?;

-- name: UpdateState :exec
UPDATE devices
SET
    state = ?
WHERE
    id = ?;

-- name: GetDevices :many
SELECT
    *
FROM
    devices;

-- name: UpdateLastSeen :exec
UPDATE devices
SET
    last_seen = ?
WHERE
    id = ?;

-- name: UpdateDeviceStateAndTrust :exec
UPDATE devices
SET
    state = ?,
    trusted_at = ?
WHERE
    id = ?;