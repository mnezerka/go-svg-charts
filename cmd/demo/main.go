package main

import (
	"log"
	"mnezerka/go-svg-charts/timestatus"
	"os"
	"text/template"
	"time"
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

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	t, err := template.New("webpage").Parse(tpl)
	check(err)

	data := struct {
		Title string
		Items []string
	}{
		Title: "go-svg-charts demo",
		Items: []string{
			timestatus.Render(
				1000,
				[]timestatus.Item{
					{Date: time.Now().Add(-24 * time.Hour), Color: "#f00"},
					{Date: time.Now().Add(-48 * time.Hour), Color: "#0f0"},
					{Date: time.Now().Add(-130 * time.Hour), Color: "#00f"},
					{Date: time.Now().Add(-140 * time.Hour), Color: "#f00"},
				},
			),
		},
	}

	err = t.Execute(os.Stdout, data)
	check(err)
}
