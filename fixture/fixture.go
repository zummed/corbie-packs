// Package fixture is what the tests use to make a project from the packs of
// this repository, as corbie add does.
package fixture

import (
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/zummed/corbie/manual"
	"github.com/zummed/corbie/pack"
)

// Root is the directory of this repository, which holds the packs.
func Root() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(file))
}

// Bot is the pack every project starts with, and Shared the pack of the
// scripts the others use.
const (
	Bot    = "bot"
	Shared = "shared"
)

// Services gives the name of every service's pack, sorted.
func Services() []string { return []string{"gitlab", "jira", "slack", "telegram"} }

// Shipped is the directory in a project that holds the packs.
const Shipped = pack.Dir

// botWord is the word a starting file holds, in its path and its text,
// where the bot's name goes.
const botWord = "BOT"

// Init adds the pack bot to a project, and the pack of each service named,
// with the bot's name, as corbie add does.
func Init(dir, bot string, services []string) error {
	words := map[string]string{}
	if bot != "" {
		words[botWord] = bot
	}
	for _, name := range append([]string{Bot}, services...) {
		src, err := pack.Parse(filepath.Join(Root(), name))
		if err != nil {
			return err
		}
		if _, err := pack.Add(dir, src, words); err != nil {
			return err
		}
	}
	return nil
}

// Starting gives the text of each starting file of every pack, by its path
// in a project, for a bot of that name.
func Starting(bot string) map[string]string {
	out := map[string]string{}
	for _, name := range append([]string{Bot}, Services()...) {
		from := filepath.Join(Root(), name, "starter")
		filepath.WalkDir(from, func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				text, _ := os.ReadFile(path)
				rel, _ := filepath.Rel(from, path)
				out[strings.ReplaceAll(filepath.ToSlash(rel), botWord, bot)] = strings.ReplaceAll(string(text), botWord, bot)
			}
			return err
		})
	}
	return out
}

// Differ gives the project's files that are not as the packs start them for
// a bot of that name, of those it has, sorted.
func Differ(dir, bot string) ([]string, error) {
	starting := Starting(bot)
	var out []string
	for _, rel := range slices.Sorted(maps.Keys(starting)) {
		mine, err := os.ReadFile(filepath.Join(dir, rel))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if string(mine) != starting[rel] {
			out = append(out, rel)
		}
	}
	return out, nil
}

// Pages are the pages the pack bot holds, each a file whose first line is
// its summary.
func Pages() []manual.Page {
	var out []manual.Page
	for _, name := range []string{"bot", "bot-project", "bot-design"} {
		text, _ := os.ReadFile(pagePath(name))
		summary, _, _ := strings.Cut(string(text), "\n")
		out = append(out, manual.Page{Name: name, Group: "Packs", Summary: summary, Text: string(text)})
	}
	return out
}

// pagePath gives the file of a page: the page bot is the pack's AGENTS.md,
// which a model is told.
func pagePath(name string) string {
	if name == "bot" {
		return filepath.Join(Root(), Bot, "AGENTS.md")
	}
	return filepath.Join(Root(), Bot, "pages", name+".txt")
}
