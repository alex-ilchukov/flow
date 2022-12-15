package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/alex-ilchukov/flow"
	"github.com/alex-ilchukov/flow/stage"
)

// A root represents root directory of the whole processing and forms paths of
// its children. It lays in the core of the Russian Doll of the flow.
type root string

func (r root) Form(j stage.Joint[int, string]) {
	w := walker{Joint: j}
	err := filepath.WalkDir(string(r), w.walk)
	if err != nil {
		j.Report(err)
	}
}

type walker struct {
	stage.Joint[int, string]
}

func (w walker) walk(path string, d fs.DirEntry, err error) error {
	if err == nil && d.Type().IsRegular() {
		w.Put(path)
	}

	return err
}

type sum struct {
	path  string
	value [md5.Size]byte
}

// md5sum is a mapper from paths to sums of files with the paths. It's an
// intermediate layer of the Russian Doll of the flow.
func md5sum(_ context.Context, path string) (sum, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return sum{}, err
	}

	return sum{path: path, value: md5.Sum(body)}, nil
}

// A collector is just a map from paths to sums of bodies of files with the
// paths. It is the outermost layer of the Russian Doll of the flow.
type collector map[string][md5.Size]byte

func (c collector) Form(j stage.Joint[sum, int]) {
	for {
		s, err := j.Get()
		if err != nil {
			return
		}

		c[s.path] = s.value
	}
}

func main() {
	flow0 := stage.Flow[int, string]{Former: root(os.Args[1])}

	flow1 := stage.Flow[string, sum]{
		Origin: &flow0,
		Former: stage.Mapper[string, sum](md5sum),
		Spread: 64,
	}

	sums := make(collector)
	flow2 := stage.Flow[sum, int]{Origin: &flow1, Former: sums}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	err := flow.Run[int](ctx, &flow2)
	if err != nil {
		cancel()
		fmt.Printf("got error: %v\n", err)
		return
	}

	paths := make([]string, 0, len(sums))
	for path := range sums {
		paths = append(paths, path)
	}

	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x  %s\n", sums[path], path)
	}
}
