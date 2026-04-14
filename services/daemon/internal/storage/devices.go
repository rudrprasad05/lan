package storage

type StoredDevice struct {
	ID        string
	Name      string
	PublicKey string
	State     string
	LastSeen  int64
	TrustedAt int64
}

func UpsertDevice(d StoredDevice) error {
	_, err := DB.Exec(`
		INSERT INTO devices (id, name, public_key, state, last_seen, trusted_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			state=excluded.state,
			last_seen=excluded.last_seen,
			trusted_at=excluded.trusted_at
	`,
		d.ID,
		d.Name,
		d.PublicKey,
		d.State,
		d.LastSeen,
		d.TrustedAt,
	)

	return err
}

func GetDevices() ([]StoredDevice, error) {
	rows, err := DB.Query(`SELECT id, name, public_key, state, last_seen, trusted_at FROM devices`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []StoredDevice

	for rows.Next() {
		var d StoredDevice
		err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.PublicKey,
			&d.State,
			&d.LastSeen,
			&d.TrustedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, d)
	}

	return list, nil
}
