package test

// The pack's lookup, placed at people/find: who someone is, by their email or
// their ID on a service, from the table people/scopes.yaml.

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

// found runs the lookup of a project with this table and input, and gives its
// output, or its error when it failed.
func found(t *testing.T, w *world, table string, input map[string]any) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "project")
	if err := exampleProject(dir, nil); err != nil {
		t.Fatal(err)
	}
	write(t, dir, map[string]string{"people/scopes.yaml": table})
	w.sync(dir)
	run, err := w.Wait(context.Background(), w.start("people/find.yaml", input), interval)
	if err != nil {
		t.Fatal(err)
	}
	if run.Status != "done" {
		return run.Status + ": " + w.show(run.ID).Error
	}
	return fmt.Sprint(w.show(run.ID).Output)
}

// The lookup finds a person by their email, ignoring case, or by a service
// and their ID there, which it reads first, and gives the table's name and
// scopes, with the email the row or the service gave. A person the table does
// not hold gets no scopes and the name the service gave, or with none their ID
// on the service. An ID with no service, or with another's, finds no row.
func TestLookupFindsAPersonByEmailOrID(t *testing.T) {
	table := `
		ana@acme.example: { name: Ana, scopes: [orders:read, orders:write] }
		bo@acme.example:  { name: Bo,  scopes: [infra:read, infra:approve] }
		"12345": { name: Bare, scopes: [infra:approve] }
		telegram:12345: { name: Tom T, scopes: [infra:read] }
		slack:U3: { scopes: [orders:read] }
		cy@acme.example: { name: Cy }
	`
	for _, c := range []struct {
		input map[string]any
		want  string
	}{
		{map[string]any{"service": "slack", "user": "U1", "name": "Bo B", "email": "Bo@Acme.Example"}, "map[email:bo@acme.example name:Bo scopes:[infra:read infra:approve]]"},
		{map[string]any{"service": "slack", "user": "U2", "name": "Eve", "email": "eve@acme.example"}, "map[email:eve@acme.example name:Eve scopes:[]]"},
		{map[string]any{"service": "telegram", "user": "12345", "name": "Tom"}, "map[email: name:Tom T scopes:[infra:read]]"},
		{map[string]any{"service": "slack", "user": "12345", "name": "Tom"}, "map[email: name:Tom scopes:[]]"},
		{map[string]any{"service": "slack", "user": "U3", "name": "Bo B", "email": "bo@acme.example"}, "map[email:bo@acme.example name:Bo B scopes:[orders:read]]"},
		{map[string]any{"service": "gitlab", "user": "7"}, "map[email: name:7 scopes:[]]"},
		{map[string]any{"email": "cy@acme.example"}, "map[email:cy@acme.example name:Cy scopes:[]]"},
	} {
		t.Run(fmt.Sprint(c.input["email"], c.input["service"], c.input["user"]), func(t *testing.T) {
			w := newWorld(t)
			if got := found(t, w, table, c.input); got != c.want {
				t.Errorf("the lookup gave %s", got)
			}
		})
	}
}

// A table with no rows finds nobody: the starter's {}, and one holding only
// comments.
func TestLookupOfAnEmptyTable(t *testing.T) {
	for _, table := range []string{"{}", "# ana@acme.example: { name: Ana, scopes: [orders:read] }"} {
		w := newWorld(t)
		if got := found(t, w, table, map[string]any{"name": "Ana", "email": "ana@acme.example"}); got != "map[email:ana@acme.example name:Ana scopes:[]]" {
			t.Errorf("with %q the lookup gave %s", table, got)
		}
	}
}

// A table that is not one fails every lookup, saying where: a line YAML cannot
// read by its number, and a row that is not { name, scopes } by its email,
// whoever is looked up.
func TestLookupRefusesAMalformedTable(t *testing.T) {
	for _, c := range []struct{ table, want string }{
		{"ana@acme.example: { name: Ana }\nbo@acme.example: name: Bo\n", "line 2"},
		{"ana@acme.example: { name: Ana }\nbo@acme.example: { name: Bo, scopes: infra:read }\n", "the row of bo@acme.example is not"},
		{"bo@acme.example: [infra:read]\n", "the row of bo@acme.example is not"},
		{"bo@acme.example: { name: Bo, scope: [infra:read] }\n", "the row of bo@acme.example is not"},
		{"- bo@acme.example\n", "people/scopes.yaml is not a mapping"},
	} {
		w := newWorld(t)
		if got := found(t, w, c.table, map[string]any{"email": "ana@acme.example"}); !strings.HasPrefix(got, "failed: ") || !strings.Contains(got, c.want) {
			t.Errorf("with %q the lookup gave %s", c.table, got)
		}
	}
}
