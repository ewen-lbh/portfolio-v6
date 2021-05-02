package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/snapcore/go-gettext"
)

func StartHTMLWatcher(messages *gettext.Catalog, db Database) {
	pugFilePattern := regexp.MustCompile(`^.+\.pug`)
	//
	// Content changes (new files or contents modified)
	//
	w := watcher.New()
	w.FilterOps(watcher.Create, watcher.Write, watcher.Move)
	w.AddFilterHook(watcher.RegexFilterHook(pugFilePattern, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				dependents := DependentsOf(event.Path, 10)
				switch event.Op {
				case watcher.Create:
					fallthrough
				case watcher.Write:
					printfln("Building file %s and its dependents %v", GetPathRelativeToSrcDir(event.Path), dependents)
					for _, filePath := range append(dependents, event.Path) {
						// Regular pages: no _ prefix
						if !strings.HasPrefix(path.Base(filePath), "_") {
							file, err := os.Stat(filePath)
							if err != nil {
								printerr("Could not stat "+filePath, err)
							}
							BuildRegularPage(messages, db, file)
						}
						if GetPathRelativeToSrcDir(filePath) == "_work.pug" {
							BuildWorkPages(db, messages)
						}
						if GetPathRelativeToSrcDir(filePath) == "_tag.pug" {
							BuildTagPages(db, messages)
						}
						if GetPathRelativeToSrcDir(filePath) == "using/_tech.pug" {
							BuildTechPages(db, messages)
						}
					}
				case watcher.Remove:
					if len(dependents) > 0 {
						printfln("WARN: Files %s depended on %s, which was removed", strings.Join(dependents, ", "), event.Path)
					}
				case watcher.Rename:
					if GetPathRelativeToSrcDir(event.OldPath) == "gallery.pug" {
						printfln("WARN: gallery.pug was renamed, exiting: you'll need to update references to the filename in Go files.")
						w.Close()
					}
					fmt.Printf("%s was renamed to %s: Updating references in %s", GetPathRelativeToSrcDir(event.OldPath), GetPathRelativeToSrcDir(event.Path), strings.Join(dependents, ", "))
					for _, filePath := range dependents {
						UpdateExtendsStatement(filePath, event.OldPath, event.Path)
					}
				}
				fmt.Println("\r\033[K")
			case err := <-w.Error:
				printerr("An errror occured while watching for content changes on *.pug files", err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive("src"); err != nil {
		printerr("Couldn't add src/ to *.pug content changes watcher", err)
	}

	if err := w.Start(100 * time.Millisecond); err != nil {
		printerr("Couldn't start the *.pug content changes watcher", err)
	}

}

// UpdateExtendsStatement renames the file referenced by an extends statement
func UpdateExtendsStatement(in string, from string, to string) {
	extendsPattern := regexp.MustCompile(`(?m)^extends\s+(?:src/)?` + from + `(?:\.pug)?\s*$`)
	file, err := os.Open(in)
	if err != nil {
		printerr(fmt.Sprintf("While updating the extends statement in %s from %s to %s: could not open file %s", in, from, to, in), err)
	}
	contents, err := os.ReadFile(in)
	if err != nil {
		printerr(fmt.Sprintf("While updating the extends statement in %s from %s to %s: could not read file %s", in, from, to, in), err)
	}
	_, err = file.Write(
		extendsPattern.ReplaceAll(contents, []byte("extends "+to)),
	)
	if err != nil {
		printerr(fmt.Sprintf("While updating the extends statement in %s from %s to %s: could not write to file %s", in, from, to, in), err)
	}
}

func GetPathRelativeToSrcDir(absPath string) string {
	return strings.SplitN(absPath, "src/", 2)[1]
}

// DependentsOf returns an array of pages' filepaths that depend
// on the given filepath (through `extends` or `intoGallery`)
// This function is recursive, dependents of dependents are also included.
// The returned array is has the same order as the build order required to correctly update dependencies before their dependents
// maxDepth is used to specify how deeply it should recurse (i.e. how many times it should call itself)
func DependentsOf(pageFilepath string, maxDepth uint) (dependents []string) {
	extendsPattern := regexp.MustCompile(fmt.Sprintf(`(?m)^extends (?:src/)?%s(?:\.pug)?$`, strings.TrimSuffix(GetPathRelativeToSrcDir(pageFilepath), ".pug")))

	err := filepath.WalkDir("src", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !dirEntry.IsDir() && strings.HasSuffix(dirEntry.Name(), ".pug") {
			content, err := os.ReadFile(path)

			if err != nil {
				printerr("Not checking for dependence on "+GetPathRelativeToSrcDir(pageFilepath)+": could not read file "+path, err)
				return nil
			}

			// If this file extends the given file,
			// or if the given file is src/gallery.pug and it uses | intoGallery (and therefore depends on src/gallery.pug)
			// add this file to the dependents
			if extendsPattern.Match(content) || (GetPathRelativeToSrcDir(pageFilepath) == "gallery.pug" && strings.Contains(string(content), "| intoGallery")) {
				dependents = append(dependents, path)
				// Add dependents of dependent after (they need to be built _after_ the dependent because they themselves depend on the former)
				if maxDepth > 1 {
					dependents = append(dependents, DependentsOf(path, maxDepth-1)...)
				} else {
					printfln("WARN: While looking for dependents for %s: Maximum recursion depth reached, not recursing any further. You might have a circular depency.", GetPathRelativeToSrcDir(pageFilepath))
				}
			}
		}
		return nil
	})

	if err != nil {
		printerr("While looking for dependents on "+GetPathRelativeToSrcDir(pageFilepath), err)
	}
	return
}
