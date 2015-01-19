package magnate

import "sort"

type migrationSorter struct {
	migrations []Migration
	by         MigrationBy
}

func (m *migrationSorter) Len() int {
	return len(m.migrations)
}

func (m *migrationSorter) Less(i, j int) bool {
	return m.by(m.migrations[i], m.migrations[j])
}

func (m *migrationSorter) Swap(i, j int) {
	m.migrations[i], m.migrations[j] = m.migrations[j], m.migrations[i]
}

type MigrationBy func(Migration, Migration) bool

func (by MigrationBy) Sort(migrations []Migration) {
	sorter := &migrationSorter{migrations, by}
	sort.Sort(sorter)
}

func Ascending(a, b Migration) bool {
	return a.Number() < a.Number()
}

func Descending(a, b Migration) bool {
	return a.Number() > b.Number()
}
