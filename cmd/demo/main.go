package main

import (
	"log"
	"os"
	"text/template"
	"time"

	"github.com/mnezerka/go-svg-charts/timestatus"
)

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		{{range .Items}}<div>{{ . }}</div>{{else}}<div><strong>no rows to be rendered</strong></div>{{end}}
	</body>
</html>`

const COL_GREEN = "#9f9"
const COL_RED = "#f99"
const COL_BLUE = "#99f"

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	t, err := template.New("webpage").Parse(tpl)
	check(err)

	cfg := timestatus.NewConfig()
	cfg.Width = 1000

	rows := []timestatus.Row{
		{
			Name: "Zero row",
			Items: []timestatus.Item{
				{Time: time.Now().Add(-24 * time.Hour), Color: COL_RED},
			},
		},
		{
			Name: "First row",
			Items: []timestatus.Item{
				{Time: time.Now().Add(-24 * time.Hour), Color: COL_RED, Label: "2-4"},
				{Time: time.Now().Add(-48 * time.Hour), Color: COL_GREEN, Label: "2-3"},
				{Time: time.Now().Add(-130 * time.Hour), Color: COL_BLUE, Label: "2-2"},
				{Time: time.Now().Add(-140 * time.Hour), Color: COL_BLUE, Label: "2-1"},
			},
		},
		{
			Name: "Second row",
			Items: []timestatus.Item{
				{Time: time.Now().Add(-24 * time.Hour), Color: COL_GREEN},
				{Time: time.Now().Add(-150 * time.Hour), Color: COL_RED},
			},
		},
	}

	data := struct {
		Title string
		Items []string
	}{
		Title: "go-svg-charts demo",
		Items: []string{timestatus.Render(rows, cfg)},
	}

	err = t.Execute(os.Stdout, data)
	check(err)
}
