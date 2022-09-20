package pdfrender

import (
	"bytes"
	"os"
	"testing"
)

func TestRenderHTML(t *testing.T) {
	r := NewPDFRenderer()

	str := `<!DOCTYPE html>
<html>
   <head>
      <title>{{.amount}}</title>
   </head>
   <body>
      <h1>{{.amount}} </h1>
      <h1>{{.name}} </h1>
      <h1>{{.address}} </h1>
      <ul >
         {{range $i, $a := .items}}
         <li>{{$a}}</li>
         {{end}}
      </ul>
      <ul >
         {{range $i, $a := .itemMap}}
         <li>{{$a}}</li>
         {{end}}
      </ul>
   </body>
</html>`

	read := bytes.NewReader([]byte(str))
	b, err := r.HTML(read)
	if err != nil {
		panic(err)
	}

	os.WriteFile("somethng.pdf", b, 0644)

}
