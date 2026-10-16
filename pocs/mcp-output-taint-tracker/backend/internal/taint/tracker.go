package taint

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

type Label string

const (
	Untrusted Label = "untrusted"
	Secret    Label = "secret"
)

const minShingleHits = 2

type TaintedEntity struct {
	Entity
	Labels []Label `json:"labels"`
}

type SinkArg struct {
	Arg    string
	Value  string
	Forbid []Label
}

type Flow struct {
	Arg        string `json:"arg"`
	Label      Label  `json:"label"`
	SourceID   int    `json:"sourceId"`
	SourceTool string `json:"sourceTool"`
	Kind       string `json:"kind"`
	EntityKind string `json:"entityKind,omitempty"`
	Evidence   string `json:"evidence"`
	Encoding   string `json:"encoding"`
	Decoded    string `json:"decoded,omitempty"`
}

type source struct {
	id     int
	tool   string
	labels []Label
}

type indexed struct {
	TaintedEntity
	lower   string
	sources []int
}

type Tracker struct {
	sources         map[int]source
	entities        map[string]*indexed
	shingles        map[uint64][]int
	trusted         []string
	trustedShingles map[uint64]bool
}

func NewTracker() *Tracker {
	return &Tracker{
		sources:         map[int]source{},
		entities:        map[string]*indexed{},
		shingles:        map[uint64][]int{},
		trustedShingles: map[uint64]bool{},
	}
}

func (t *Tracker) AddTrusted(text string) {
	t.trusted = append(t.trusted, strings.ToLower(text))
	for _, h := range shingles(tokenize(text)) {
		t.trustedShingles[h] = true
	}
}

func (t *Tracker) AddSource(id int, tool string, labels []Label, text string) []TaintedEntity {
	t.sources[id] = source{id, tool, labels}
	var out []TaintedEntity
	for _, f := range extract(text) {
		ls := slices.Clone(labels)
		if f.secret && !slices.Contains(ls, Secret) {
			ls = append(ls, Secret)
		}
		if len(ls) == 0 {
			continue
		}
		key := f.Kind + "|" + strings.ToLower(f.Value)
		e, ok := t.entities[key]
		if !ok {
			e = &indexed{TaintedEntity: TaintedEntity{Entity: f.Entity}, lower: strings.ToLower(f.Value)}
			t.entities[key] = e
		}
		e.Labels = union(e.Labels, ls)
		e.sources = append(e.sources, id)
		out = append(out, TaintedEntity{f.Entity, ls})
	}
	if len(labels) > 0 {
		for _, h := range shingles(tokenize(text)) {
			if !slices.Contains(t.shingles[h], id) {
				t.shingles[h] = append(t.shingles[h], id)
			}
		}
	}
	return out
}

func (t *Tracker) Check(args []SinkArg) []Flow {
	type pick struct {
		flow Flow
		size int
	}
	best := map[string]pick{}
	var order []string
	keep := func(f Flow, matched string) {
		key := fmt.Sprintf("%s|%s|%s|%d", f.Arg, f.Label, f.Kind, f.SourceID)
		prev, ok := best[key]
		if !ok {
			order = append(order, key)
		}
		if !ok || len(matched) > prev.size {
			best[key] = pick{f, len(matched)}
		}
	}
	for _, a := range args {
		for _, v := range variants(a.Value) {
			t.matchEntities(a, v, keep)
			t.matchText(a, v, keep)
		}
	}
	out := make([]Flow, 0, len(order))
	for _, k := range order {
		out = append(out, best[k].flow)
	}
	return out
}

func (t *Tracker) matchEntities(a SinkArg, v variant, keep func(Flow, string)) {
	low := strings.ToLower(v.text)
	for _, e := range t.sortedEntities() {
		if !strings.Contains(low, e.lower) {
			continue
		}
		for _, l := range e.Labels {
			if !slices.Contains(a.Forbid, l) || (l == Untrusted && t.endorsed(e.lower)) {
				continue
			}
			for _, id := range e.sources {
				src := t.sources[id]
				if !slices.Contains(src.labels, l) && !(l == Secret && e.Kind == "secret") {
					continue
				}
				keep(t.flow(a, l, src, "entity", e.Kind, v, e.Value), e.Value)
			}
		}
	}
}

func (t *Tracker) matchText(a SinkArg, v variant, keep func(Flow, string)) {
	toks := tokenize(v.text)
	hits := map[int][]int{}
	for i, h := range shingles(toks) {
		if t.trustedShingles[h] {
			continue
		}
		for _, id := range t.shingles[h] {
			hits[id] = append(hits[id], i)
		}
	}
	for id, pos := range hits {
		if len(pos) < minShingleHits {
			continue
		}
		src := t.sources[id]
		start := toks[pos[0]].start
		end := toks[pos[len(pos)-1]+shingleSize-1].end
		matched := v.text[start:end]
		for _, l := range src.labels {
			if slices.Contains(a.Forbid, l) {
				keep(t.flow(a, l, src, "text", "", v, matched), matched)
			}
		}
	}
}

func (t *Tracker) flow(a SinkArg, l Label, src source, kind, entityKind string, v variant, matched string) Flow {
	f := Flow{Arg: a.Arg, Label: l, SourceID: src.id, SourceTool: src.tool, Kind: kind, EntityKind: entityKind, Encoding: v.encoding}
	if v.encoding == "plain" {
		f.Evidence = matched
	} else {
		f.Evidence = v.raw
		f.Decoded = matched
	}
	return f
}

func (t *Tracker) endorsed(lower string) bool {
	for _, s := range t.trusted {
		if strings.Contains(s, lower) {
			return true
		}
	}
	return false
}

func (t *Tracker) sortedEntities() []*indexed {
	out := make([]*indexed, 0, len(t.entities))
	for _, e := range t.entities {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].lower < out[j].lower })
	return out
}

func union(a, b []Label) []Label {
	out := slices.Clone(a)
	for _, l := range b {
		if !slices.Contains(out, l) {
			out = append(out, l)
		}
	}
	return out
}
