package timestatus

import (
	"fmt"
	"os"
	"sort"
	"time"
)

const BORDER = 2
const ROW_LABEL_WIDTH = 0

type Row struct {
	Name  string
	Items []Item
}

type Item struct {
	Time  time.Time
	Color string
	Label string
}

const ITEM_HEIGHT = 20

func render_row(width int, start time.Time, end time.Time, row Row, offset_y int) string {
	result := ""

	fmt.Fprintf(os.Stderr, "Rendering row \"%s\", %d items, start=%s, end=%s\n", row.Name, len(row.Items), start, end)

	if len(row.Items) == 0 {
		return result
	}

	//var first_item Item
	var last_offset = 0
	last := start

	// width ot the line in unix nanoseconds
	range_unix := end.Unix() - start.Unix()

	// don't render first item, just remember its date
	for i := 0; i < len(row.Items); i++ {

		item := row.Items[i]

		fmt.Fprintf(os.Stderr, "%s --- %s --- %s\n", row.Items[i], start, end)

		item_width := float64(item.Time.Unix()-last.Unix()) / float64(range_unix) * float64(width-ROW_LABEL_WIDTH)

		fmt.Fprintf(os.Stderr, "width=%f\n", item_width)

		//fmt.Printf("width_ratio: %f, width: %d\n", item_width_ratio, int(item_width))

		// item rect
		result += fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" style="fill:%s;stroke: none;" />`, int(last_offset), offset_y, int(item_width)-BORDER, ITEM_HEIGHT-BORDER, item.Color)
		// item circle
		//result += fmt.Sprintf(`<circle cx="%d" cy="%d" r="%d" style="fill: black ;stroke: gray;" />`, int(last_offset)+int(item_width), offset_y+int(ITEM_HEIGHT/2.0), 4)
		// item label
		if len(item.Label) > 0 {
			result += fmt.Sprintf(`<text x="%d" y="%d", style="color: white; font-size: 13px;">%s</text>`, int(last_offset)+2, offset_y+ITEM_HEIGHT-4, item.Label)
		}

		// new last offset
		last_offset += int(item_width)
	}

	return result
}

func Render(width int, rows []Row) string {

	result := ""

	if len(rows) == 0 {
		return result
	}

	start := time.Now()
	end := start

	// sort items in individual rows and compute the oldest item
	for i := 0; i < len(rows); i++ {
		row := rows[i]

		sort.Slice(row.Items, func(i, j int) bool {
			return row.Items[i].Time.Before(row.Items[j].Time)
		})

		if len(row.Items) > 0 {
			if row.Items[0].Time.Before(start) {
				start = row.Items[0].Time
			}
		}
	}

	total_range := end.Sub(start)
	// let's add 1/10 of total range at the beginning (before first item)
	start = start.Add(-1 * total_range / 10)

	result += fmt.Sprintf(`<svg width="%d" height="110">`, width)

	for i := 0; i < len(rows); i++ {
		result += render_row(width, start, end, rows[i], i*ITEM_HEIGHT)
	}

	// row rectangle
	result += fmt.Sprintf(`<rect width="%d" height="%d" style="fill:none;stroke-width:1;stroke: #000 " />`, width, ITEM_HEIGHT*len(rows))

	result += `</svg>`

	return result
}
