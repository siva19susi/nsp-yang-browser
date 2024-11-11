package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/openconfig/goyang/pkg/yang"
)

func (a *App) readYangModules(files []string) error {
	if len(files) == 0 {
		return nil
	}

	for _, name := range files {
		if err := a.modules.Read(name); err != nil {
			return err
		}
	}

	if errors := a.modules.Process(); len(errors) > 0 {
		errorStrings := make([]string, len(errors))
		for i, e := range errors {
			errorStrings[i] = e.Error()
		}
		return fmt.Errorf("yang processing failed: %s", strings.Join(errorStrings, ", "))
	}

	return a.generateYangSchema()
}

func (a *App) readNspYangModules(modules []IntentTypeYangModule) error {
	for _, entry := range modules {
		if err := a.modules.Parse(entry.YangContent, entry.Name); err != nil {
			return err
		}
	}

	if errors := a.modules.Process(); len(errors) > 0 {
		errorStrings := make([]string, len(errors))
		for i, e := range errors {
			errorStrings[i] = e.Error()
		}
		return fmt.Errorf("yang processing failed: %s", strings.Join(errorStrings, ", "))
	}

	return a.generateYangSchema()
}

func (a *App) generateYangSchema() error {
	// Keep track of the top level modules we read in.
	// Those are the only modules we want to print below.
	mods := map[string]*yang.Module{}
	var names []string

	for _, m := range a.modules.Modules {
		if mods[m.Name] == nil {
			mods[m.Name] = m
			names = append(names, m.Name)
		}
	}
	sort.Strings(names)
	entries := make([]*yang.Entry, len(names))
	for x, n := range names {
		entries[x] = yang.ToEntry(mods[n])
	}

	a.SchemaTree = buildRootEntry()

	for _, entry := range entries {
		updateAnnotation(entry)
		a.SchemaTree.Dir[entry.Name] = entry
	}
	return nil
}

func buildRootEntry() *yang.Entry {
	return &yang.Entry{
		Name: "root",
		Kind: yang.DirectoryEntry,
		Dir:  make(map[string]*yang.Entry),
		Annotation: map[string]interface{}{
			"schemapath": "/",
			"root":       true,
		},
	}
}

// updateAnnotation updates the schema info before encoding.
func updateAnnotation(entry *yang.Entry) {
	for _, child := range entry.Dir {
		updateAnnotation(child)
		child.Annotation = map[string]interface{}{}
		t := child.Type
		if t == nil {
			continue
		}

		switch t.Kind {
		case yang.Ybits:
			nameMap := t.Bit.NameMap()
			bits := make([]string, 0, len(nameMap))
			for bitstr := range nameMap {
				bits = append(bits, bitstr)
			}
			child.Annotation["bits"] = bits
		case yang.Yenum:
			nameMap := t.Enum.NameMap()
			enum := make([]string, 0, len(nameMap))
			for enumstr := range nameMap {
				enum = append(enum, enumstr)
			}
			child.Annotation["enum"] = enum
		case yang.Yidentityref:
			identities := make([]string, 0, len(t.IdentityBase.Values))
			for i := range t.IdentityBase.Values {
				identities = append(identities, t.IdentityBase.Values[i].PrefixedName())
			}
			child.Annotation["prefix-qualified-identities"] = identities
		}
		if t.Root != nil {
			child.Annotation["root.type"] = t.Root.Name
		}
	}
}
