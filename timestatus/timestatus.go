package timestatus

import (
	"fmt"
	"os"
	"sort"
	"time"
)

type Config struct {
	Width             int
	RowLabelWidth     int
	RenderBoundingBox bool
	RenderItemLabels  bool
	ItemHeight        int
	ItemBorder        int
	AxisXParts        int
}

func NewConfig() Config {
	config := Config{}
	config.Width = 1000
	config.RowLabelWidth = 100
	config.RenderBoundingBox = true
	config.RenderItemLabels = true
	config.ItemHeight = 20
	config.ItemBorder = 2
	config.AxisXParts = 5

	return config
}

type Row struct {
	Name  string
	Items []Item
}

type Item struct {
	Time  time.Time
	Color string
	Label string
	URL   string
}

func raw_time(t time.Time) string {
	return t.Format("2006/01/02 03:04")
}

func format_date_short(t time.Time) string {
	return t.Format("01/02")
}

func render_axis_x(config Config, start time.Time, end time.Time, offset_y int) string {

	var result = ""

	fmt.Fprintf(os.Stderr, "Rendering axis x\n")

	dur := end.Sub(start)
	dur_part := dur / 10
	dur_part_width := float64(config.Width-config.RowLabelWidth) / float64(config.AxisXParts)

	fmt.Fprintf(os.Stderr, "  dur: %v\n", dur)
	fmt.Fprintf(os.Stderr, "  part: %v\n", dur_part)
	fmt.Fprintf(os.Stderr, "  part width: %f\n", dur_part_width)
	fmt.Fprintf(os.Stderr, "  hours: %v\n", dur_part.Hours())
	fmt.Fprintf(os.Stderr, "  hours: %v\n", dur_part.Truncate(time.Hour))

	t := start

	for i := 0; i < config.AxisXParts; i++ {
		pos_x := config.RowLabelWidth + int(float64(i)*dur_part_width+dur_part_width/2.0)
		fmt.Fprintf(os.Stderr, "  iter: %d, time: %s, pos: %d\n", i, raw_time(t), pos_x)

		result += fmt.Sprintf(
			`<text x="%d" y="%d" dominant-baseline="middle" text-anchor="middle" style="color: white; font-size: 13px;">%s</text>`,
			pos_x,
			offset_y+int(config.ItemHeight/2.0),
			format_date_short(t))

		t = t.Add(dur_part)
	}

	return result
}

func render_row(config Config, start time.Time, end time.Time, row Row, offset_y int) string {
	result := ""

	fmt.Fprintf(os.Stderr, "Rendering row \"%s\", %d items, start=%s, end=%s\n", row.Name, len(row.Items), raw_time(start), raw_time(end))

	if len(row.Items) == 0 {
		return result
	}

	result += fmt.Sprintf(
		`<text x="0" y="%d" dominant-baseline="middle" style="color: black; font-size: 13px;">%s</text>`,
		offset_y+int(config.ItemHeight/2.0),
		row.Name)

	//var first_item Item
	var last_offset = 0
	last := start

	// width ot the line in unix nanoseconds
	range_unix := end.Unix() - start.Unix()

	// don't render first item, just remember its date
	for i := 0; i < len(row.Items); i++ {

		item := row.Items[i]

		fmt.Fprintf(os.Stderr, "  rendering item %s, duration: %v", item.Label, item.Time.Sub(last))

		item_width := float64(item.Time.Unix()-last.Unix()) / float64(range_unix) * float64(config.Width-config.RowLabelWidth)

		fmt.Fprintf(os.Stderr, "    width=%f\n", item_width)

		// item rect
		//result += fmt.Sprintf(`<a href="https://seznam.cz">`)
		result += fmt.Sprintf(
			`<rect x="%d" y="%d" width="%d" height="%d" style="fill:%s;stroke: none;" />`,
			int(last_offset)+config.RowLabelWidth,
			offset_y,
			int(item_width)-config.ItemBorder,
			config.ItemHeight-config.ItemBorder,
			item.Color)
		//result += fmt.Sprintf(`</a>`)
		// item label
		if len(item.Label) > 0 {
			result += fmt.Sprintf(
				`<text x="%d" y="%d" dominant-baseline="middle" style="color: white; font-size: 13px;">%s</text>`,
				int(last_offset)+2+config.RowLabelWidth,
				offset_y+int(config.ItemHeight/2.0),
				item.Label)
		}

		// new last offset
		last_offset += int(item_width)
		last = item.Time
	}

	return result
}

func Render(rows []Row, config Config) string {

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

	result += fmt.Sprintf(`<svg width="%d" height="110">`, config.Width)

	for i := 0; i < len(rows); i++ {
		result += render_row(config, start, end, rows[i], i*config.ItemHeight)
	}

	// row rectangle
	result += fmt.Sprintf(`<rect width="%d" height="%d" style="fill:none;stroke-width:1;stroke: #000 " />`, config.Width, config.ItemHeight*len(rows))

	result += render_axis_x(config, start, end, len(rows)*config.ItemHeight)

	result += `</svg>`

	return result
}
