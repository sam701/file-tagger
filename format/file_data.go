package format

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aybabtme/rgbterm"
)

type FileData struct {
	Id           uint64
	Name         string
	Period       string
	Tags         []string
	CreationDate time.Time
}

type Files []*FileData

func (fs Files) Print() {
	sort.Sort(byPeriod(fs))

	prevPeriod := ""
	for _, f := range fs {
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

}

type byPeriod []*FileData

func (a byPeriod) Len() int           { return len(a) }
func (a byPeriod) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPeriod) Less(i, j int) bool { return a[i].Period > a[j].Period }
