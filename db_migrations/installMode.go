package dbmigrations

var installModequeries = []string{
	"CREATE TABLE IF NOT EXISTS db_version_history(name VARCHAR(50), version INT, description VARCHAR(200), updated_at BIGINT)",
	"CREATE TABLE IF NOT EXISTS db_version (version INT)",
	"INSERT INTO db_version (version) VALUES (0)",
}
