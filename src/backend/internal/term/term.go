// Package term contains the academic term summary shared by feature packages.
package term

// Term is an academic term. Dates are "YYYY-MM-DD" in campus local time.
type Term struct {
	ID       int64  `json:"id,string"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	StartsOn string `json:"starts_on"`
	EndsOn   string `json:"ends_on"`
}

// Columns selects a Term from the terms table aliased as t. Scan it with Dest.
const Columns = "t.id, t.code, t.name, to_char(t.starts_on, 'YYYY-MM-DD'), to_char(t.ends_on, 'YYYY-MM-DD')"

// Dest returns the scan destinations matching Columns.
func Dest(t *Term) []any {
	return []any{&t.ID, &t.Code, &t.Name, &t.StartsOn, &t.EndsOn}
}
