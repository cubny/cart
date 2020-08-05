package sqlite3

import "fmt"

// Migrate migrates the database
// there are much better tools for performing migrations, but for
// the purpose of the assignment, we keep things as simple as possible
func (s *Sqlite3) Migrate() error {
	migrations := []string{
		migration01MigrationCreateCartsTable,
		migration02AddCartIndex,
		migration03CreateLineItemsTable,
		migration04AddLineItemsIndex,
		migration05AddLineItemsIndex2,
	}

	for i, m := range migrations {
		if _, err := s.db.Exec(m); err != nil {
			return fmt.Errorf("migration failed at index of %d, %s", i, err)
		}
	}

	return nil
}

// truncate truncates all tables
// it is meant to be used for integration tests
func (s *Sqlite3) TruncateAllTables() error {
	truncates := []string{
		truncateCartsTable,
		truncateLineItemsTable,
	}

	for i, m := range truncates {
		if _, err := s.db.Exec(m); err != nil {
			return fmt.Errorf("truncate failed at index of %d, %s", i, err)
		}
	}
	return nil
}
