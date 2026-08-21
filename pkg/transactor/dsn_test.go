package transactor

import (
	"database/sql"
	"path/filepath"
	"testing"
)

// TestDSNAppliesTheProfile opens a real file and asks the database what it got.
// The DSN string is exactly what differs between the two backends — each ignores
// the other's spelling — so asserting on the string would pass on the build that
// ignores it.
func TestDSNAppliesTheProfile(t *testing.T) {
	db, err := sql.Open("sqlite", dsn(filepath.Join(t.TempDir(), "t.db"), "WAL"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = db.Close() }()

	for _, c := range []struct{ pragma, want string }{
		{"journal_mode", "wal"},
		{"busy_timeout", "5000"},
		{"foreign_keys", "0"}, // deliberately off: docs has no references
	} {
		var got string
		if err := db.QueryRow("PRAGMA " + c.pragma).Scan(&got); err != nil {
			t.Fatalf("PRAGMA %s: %v", c.pragma, err)
		}
		if got != c.want {
			t.Errorf("%s = %q, want %q", c.pragma, got, c.want)
		}
	}
}
