package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/openconfig/goyang/pkg/yang"

	"github.com/openconfig/gnmic/pkg/api/path"
)

type App struct {
	SchemaTree *yang.Entry
	modules    *yang.Modules
}

type generatedPath struct {
	Path           string   `json:"path,omitempty"`
	PathWithPrefix string   `json:"path-with-prefix,omitempty"`
	Type           string   `json:"type,omitempty"`
	EnumValues     []string `json:"enum-values,omitempty"`
	Description    string   `json:"description,omitempty"`
	Default        string   `json:"default,omitempty"`
	IsState        bool     `json:"is-state,omitempty"`
	Namespace      string   `json:"namespace,omitempty"`
	FeatureList    []string `json:"if-features,omitempty"`
	IsNotification bool     `json:"is-notification,omitempty"`
	IsRpc          bool     `json:"is-rpc,omitempty"`
	IsAction       bool     `json:"is-action,omitempty"`
}

func (a *App) pathCmdRun() ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := make(chan *generatedPath)
	gpaths := make([]*generatedPath, 0)
	done := make(chan struct{})
	go func(ctx context.Context, out chan *generatedPath) {
		for {
			select {
			case m, ok := <-out:
				if !ok {
					close(done)
					return
				}
				gpaths = append(gpaths, m)
			case <-ctx.Done():
				return
			}
		}
	}(ctx, out)

	collected := make([]*yang.Entry, 0, 256)

	withNonLeaves := true
	for _, entry := range a.SchemaTree.Dir {
		collected = append(collected, collectSchemaNodes(entry, withNonLeaves)...)
	}

	for _, entry := range collected {
		// don't produce such paths in case of non-leaves
		if entry.IsCase() || entry.IsChoice() {
			continue
		}
		out <- generatePath(entry)
	}
	close(out)
	<-done

	sort.Slice(gpaths, func(i, j int) bool {
		return gpaths[i].Path < gpaths[j].Path
	})
	for _, gp := range gpaths {
		gp.PathWithPrefix = collapsePrefixes(gp.PathWithPrefix)
	}

	if len(gpaths) == 0 {
		return []byte("No results found"), nil
	}

	b, err := json.MarshalIndent(gpaths, "", "  ")
	if err != nil {
		return []byte{}, fmt.Errorf("[Error] generating paths: %s", err)
	}

	return b, nil
}

func collectSchemaNodes(e *yang.Entry, leafOnly bool) []*yang.Entry {
	if e == nil {
		return []*yang.Entry{}
	}
	collected := make([]*yang.Entry, 0, 128)
	for _, child := range e.Dir {
		collected = append(collected,
			collectSchemaNodes(child, leafOnly)...)
	}

	// Support for Notification
	if e.Node.Kind() == "notification" {
		if e.Extra == nil {
			e.Extra = make(map[string][]interface{})
		}
		e.Extra["notification"] = []interface{}{true}
	}

	// Support for RPC & Action
	if e.RPC != nil {
		kind := e.Node.Kind()
		if e.RPC.Input != nil {
			if e.RPC.Input.Extra == nil {
				e.RPC.Input.Extra = make(map[string][]interface{})
			}
			e.RPC.Input.Extra[kind] = []interface{}{true}
			collected = append(collected, collectSchemaNodes(e.RPC.Input, leafOnly)...)
		}
		if e.RPC.Output != nil {
			if e.RPC.Output.Extra == nil {
				e.RPC.Output.Extra = make(map[string][]interface{})
			}
			e.RPC.Output.Extra[kind] = []interface{}{true}
			collected = append(collected, collectSchemaNodes(e.RPC.Output, leafOnly)...)
		}
	}

	if e.Parent != nil {
		switch {
		case e.Dir == nil && e.ListAttr != nil: // leaf-list
			fallthrough
		case e.Dir == nil: // leaf
			f := &yang.Entry{
				Parent:      e.Parent,
				Node:        e.Node,
				Name:        e.Name,
				Description: e.Description,
				Default:     e.Default,
				Units:       e.Units,
				Kind:        e.Kind,
				Config:      e.Config,
				Prefix:      e.Prefix,
				Mandatory:   e.Mandatory,
				Dir:         e.Dir,
				Key:         e.Key,
				Type:        e.Type,
				Exts:        e.Exts,
				ListAttr:    e.ListAttr,
				Extra:       make(map[string][]any),
			}
			for k, v := range e.Extra {
				f.Extra[k] = v
			}
			collected = append(collected, f)
		case e.ListAttr != nil: // list
			fallthrough
		default: // container
			if !leafOnly {
				collected = append(collected, e)
			}
			if len(e.Extra["if-feature"]) > 0 {
				for _, myleaf := range collected {
					if myleaf.Extra["if-feature"] == nil {
						myleaf.Extra["if-feature"] = e.Extra["if-feature"]
						continue
					}
				LOOP:
					for _, f := range e.Extra["if-feature"] {
						for _, mlf := range myleaf.Extra["if-feature"] {
							if ff, ok := f.(*yang.Value); ok && ff != nil {
								if mlff, ok := mlf.(*yang.Value); ok && mlff != nil {
									if ff.Source == nil || mlff.Source == nil {
										continue LOOP
									}
									if ff.Source.Argument == mlff.Source.Argument {
										continue LOOP
									}
									myleaf.Extra["if-feature"] = append(myleaf.Extra["if-feature"], f)
								}
							}
						}
					}
				}
			}
			// Support for Notification
			if len(e.Extra["notification"]) > 0 {
				for _, myleaf := range collected {
					if myleaf.Extra["notification"] == nil {
						myleaf.Extra["notification"] = e.Extra["notification"]
					}
				}
			}
			// Support for RPC
			if len(e.Extra["rpc"]) > 0 {
				for _, myleaf := range collected {
					if myleaf.Extra["rpc"] == nil {
						myleaf.Extra["rpc"] = e.Extra["rpc"]
					}
				}
			}
			// Support for Action
			if len(e.Extra["action"]) > 0 {
				for _, myleaf := range collected {
					if myleaf.Extra["action"] == nil {
						myleaf.Extra["action"] = e.Extra["action"]
					}
				}
			}
		}
	}
	return collected
}

func generatePath(entry *yang.Entry) *generatedPath {
	gp := new(generatedPath)
	for e := entry; e != nil && e.Parent != nil; e = e.Parent {
		if e.IsCase() || e.IsChoice() {
			continue
		}
		elementName := e.Name
		prefixedElementName := e.Name
		if e.Prefix != nil {
			prefixedElementName = fmt.Sprintf("%s:%s", e.Prefix.Name, prefixedElementName)
		}
		if e.Key != "" {
			for _, k := range strings.Fields(e.Key) {
				elementName = fmt.Sprintf("%s[%s=*]", elementName, k)
				prefixedElementName = fmt.Sprintf("%s[%s=*]", prefixedElementName, k)
			}
		}
		gp.Path = fmt.Sprintf("/%s%s", elementName, gp.Path)
		if e.Prefix != nil {
			gp.PathWithPrefix = fmt.Sprintf("/%s%s", prefixedElementName, gp.PathWithPrefix)
		}
	}
	if ifFeature, ok := entry.Extra["if-feature"]; ok && ifFeature != nil {
	APPEND:
		for _, feature := range ifFeature {
			f, ok := feature.(*yang.Value)
			if !ok {
				continue
			}
			for _, ef := range gp.FeatureList {
				if ef == f.Source.Argument {
					continue APPEND
				}
			}
			gp.FeatureList = append(gp.FeatureList, strings.Split(f.Source.Argument, " and ")...)
		}
	}

	// Support for Notification
	if len(entry.Extra["notification"]) == 1 {
		gp.IsNotification = true
	}
	// Support for RPC
	if len(entry.Extra["rpc"]) == 1 {
		gp.IsRpc = true
	}
	// Support for Action
	if len(entry.Extra["action"]) == 1 {
		gp.IsAction = true
	}

	gp.Description = entry.Description
	if entry.Type != nil {
		gp.Type = entry.Type.Name
		if gp.Type == "enumeration" {
			gp.EnumValues = entry.Type.Enum.Names()
		}
	} else if entry.IsList() {
		gp.Type = "[list]"
	} else {
		gp.Type = "[container]"
	}

	if entry.IsLeafList() {
		gp.Default = strings.Join(entry.DefaultValues(), ", ")
	} else {
		gp.Default, _ = entry.SingleDefaultValue()
	}

	gp.IsState = isState(entry)
	gp.Namespace = entry.Namespace().NName()
	return gp
}

func isState(e *yang.Entry) bool {
	if e.Config == yang.TSFalse {
		return true
	}
	if e.Parent != nil {
		return isState(e.Parent)
	}
	return false
}

// collapsePrefixes removes prefixes from path element names and keys
func collapsePrefixes(p string) string {
	gp, err := path.ParsePath(p)
	if err != nil {
		return p
	}
	parentPrefix := ""
	for _, pe := range gp.Elem {
		currentPrefix, name := getPrefixElem(pe.Name)
		if parentPrefix == "" || parentPrefix != currentPrefix {
			// first elem or updating parent prefix
			parentPrefix = currentPrefix
		} else if currentPrefix == parentPrefix {
			pe.Name = name
		}
	}
	return fmt.Sprintf("/%s", path.GnmiPathToXPath(gp, false))
}

// takes a path element name or a key name
// and returns the prefix and name
func getPrefixElem(pe string) (string, string) {
	if pe == "" {
		return "", ""
	}
	pes := strings.SplitN(pe, ":", 2)
	if len(pes) > 1 {
		return pes[0], pes[1]
	}
	return "", pes[0]
}
