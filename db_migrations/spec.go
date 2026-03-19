package dbmigrations

type DBVersionUpgrade struct {
	Name        string `db:"name"`
	Version     int    `db:"version"`
	Description string `db:"description"`
	UpdatedAt   int64  `db:"updated_at"`
}

type DBVersion struct {
	Version int `db:"version"`
}

const (
	Name        = "name"
	Version     = "version"
	Description = "description"
	UpdatedAT   = "updated_at"
)

var DBVersionHistInsertColumns = []string{Name, Version, Description, UpdatedAT}

const (
	DBVersionHistoryTableName = "db_version_history"
)
