package codenames

import (
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
)

const tpl = `
<!DOCTYPE html>
<html>
    <head>
        <title>Codenames - Play Online</title>
        <script src="/static/app.js?v=0.02" type="text/javascript"></script>
        <link href="https://fonts.googleapis.com/css?family=Roboto" rel="stylesheet">
        <link rel="stylesheet" type="text/css" href="/static/game.css" />
        <link rel="stylesheet" type="text/css" href="/static/lobby.css" />
        <link rel="shortcut icon" type="image/png" id="favicon" href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAA8SURBVHgB7dHBDQAgCAPA1oVkBWdzPR84kW4AD0LCg36bXJqUcLL2eVY/EEwDFQBeEfPnqUpkLmigAvABK38Grs5TfaMAAAAASUVORK5CYII="/>

        <script type="text/javascript">
             {{if .SelectedGameID}}
             window.selectedGameID = "{{.SelectedGameID}}";
             {{end}}
             window.autogeneratedGameID = "{{.AutogeneratedGameID}}";
        </script>
    </head>
    <body>
		<script>

		  ga('create', 'UA-88084599-2', 'auto');
		  ga('send', 'pageview');

		</script>
		<div id="app">
		</div>
    </body>
</html>
`

type templateParameters struct {
	SelectedGameID      string
	AutogeneratedGameID string
}

func (s *Server) handleIndex(rw http.ResponseWriter, req *http.Request) {
	dir, id := filepath.Split(req.URL.Path)
	if dir != "" && dir != "/" {
		http.NotFound(rw, req)
		return
	}

	autogeneratedID := s.getAutogeneratedID()

	err := s.tpl.Execute(rw, templateParameters{
		SelectedGameID:      id,
		AutogeneratedGameID: autogeneratedID,
	})
	if err != nil {
		http.Error(rw, "error rendering", http.StatusInternalServerError)
	}
}

func (s *Server) getAutogeneratedID() string {
	const attemptsPerWordCount = 5

	s.mu.Lock()
	defer s.mu.Unlock()

	var words []string
	autogeneratedID := ""
	for i := 0; ; i++ {
		wordCount := 2 + i/attemptsPerWordCount

		words = words[:0]
		for j := 0; j < wordCount; j++ {
			w := s.gameIDWords[rand.Intn(len(s.gameIDWords))]
			words = append(words, w)
		}

		autogeneratedID = strings.Join(words, "-")
		if _, ok := s.games[autogeneratedID]; !ok {
			break
		}
	}
	return autogeneratedID
}
