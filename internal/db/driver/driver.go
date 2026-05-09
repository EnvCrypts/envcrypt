package driver

import "fmt"

const (
	AppDriverPostgres = "postgres"
	AppDriverSQLite   = "sqlite"
)

func SQLDriverName(appDriver string) (string, error) {
	switch normalize(appDriver) {
	case AppDriverPostgres:
		return "pgx", nil
	case AppDriverSQLite:
		return "sqlite", nil
	default:
		return "", fmt.Errorf("unsupported database driver %q", appDriver)
	}
}

func IsSQLite(appDriver string) bool {
	return normalize(appDriver) == AppDriverSQLite
}

func normalize(appDriver string) string {
	switch appDriver {
	case "", "postgresql":
		return AppDriverPostgres
	default:
		return appDriver
	}
}
