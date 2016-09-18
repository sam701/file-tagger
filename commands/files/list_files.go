package files

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aybabtme/rgbterm"
	"github.com/sam701/file-tagger/storage"
	"github.com/urfave/cli"
)

func List(c *cli.Context) error {
	files := storage.Open(c).GetFiles(c.Args())
	sort.Sort(byPeriod(files))

	prevPeriod := ""
	for _, f := range files {
		if prevPeriod != f.Period {
			fmt.Println(rgbterm.FgString(f.Period, 255, 255, 255))
			prevPeriod = f.Period
		}
		id := rgbterm.FgString(fmt.Sprintf("%08x", f.Id), 100, 255, 100)
		allTags := strings.Join(f.Tags, " ")
		allTags = fmt.Sprintf("%-15s", allTags)
		tags := rgbterm.FgString(allTags, 255, 100, 100)

		fmt.Println(" ", id, tags, f.Name)
	}
	return nil
}

type byPeriod []*storage.FileData

func (a byPeriod) Len() int           { return len(a) }
func (a byPeriod) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPeriod) Less(i, j int) bool { return a[i].Period > a[j].Period }
