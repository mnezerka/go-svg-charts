package timestatus

import (
	"fmt"
	"sort"
	"time"
)

type Item struct {
	Date  time.Time
	Color string
}

const ITEM_HEIGHT = 20

func Render(width int, items []Item) string {
	result := ""

	sort.Slice(items, func(i, j int) bool {
		return items[i].Date.Before(items[j].Date)
	})

	result += fmt.Sprintf(`<svg width="%d" height="110">`, width)

	//var first_item Item
	var range_unix int64 = 0
	var start_unix int64 = 0
	var last_unix int64 = 0
	var last_offset = 0
	for i := 0; i < len(items); i++ {
		item := items[i]

		// remember first date
		if i == 0 {
			//first_item = item
			start_unix = item.Date.Unix()
			range_unix = time.Now().Unix() - item.Date.Unix()
			last_unix = start_unix
		}

		// this is preparation for rendering multiple rows, where first item
		// in a row doesn't have to aligned with left border of canvas
		if item.Date.Unix() == start_unix {
			continue
		}

		item_width := float64(item.Date.Unix()-last_unix) / float64(range_unix) * float64(width)

		//fmt.Printf("width_ratio: %f, width: %d\n", item_width_ratio, int(item_width))

		//result += fmt.Sprintf(`<rect width="%s" height="100" style="fill:rgb(0,0,255);stroke-width:1;stroke: " />`, item_width, item.Color)
		result += fmt.Sprintf(`<rect x="%d" y="0" width="%d" height="%d" style="fill:%s;stroke: none;" />`, int(last_offset), int(item_width), ITEM_HEIGHT, item.Color)
		result += fmt.Sprintf(`<text x="%d" y="%d", style="color: white; font-size: 13px;">BUILD</text>`, int(last_offset)+2, ITEM_HEIGHT-4)

		// new last offset
		last_offset += int(item_width)
	}

	// row rectangle
	result += fmt.Sprintf(`<rect width="%d" height="%d" style="fill:none;stroke-width:1;stroke: #000 " />`, width, ITEM_HEIGHT)

	result += `</svg>`

	return result
}
