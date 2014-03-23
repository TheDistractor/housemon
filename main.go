package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jcw/flow"
	_ "github.com/jcw/jeebus/gadgets"
)

var settings = map[string]string{
	"appDir": "./app",
	"baseDir": "./base",
	"dataDir": "./data",
	"port": "5561",
}

func main() {
	loadSettings("./settings.txt")
	
	fmt.Printf("Starting webserver at http://localhost:%s/\n", settings["port"])
	
	// show an intro page via a special webserver if the main app dir is absent
	if !fileExists(settings["appDir"]) {
		startIntroServer() // never returns
	}

	c := flow.NewCircuit()
	c.Add("http", "HTTPServer")
	c.Add("forever", "Forever")
	c.Feed("http.Handlers", flow.Tag{"/", settings["appDir"]})
	c.Feed("http.Handlers", flow.Tag{"/base", settings["baseDir"]})
	c.Feed("http.Handlers", flow.Tag{"/ws", "<websocket>"})
	c.Feed("http.Start", settings["port"])
	c.Run()
}

// startIntroServer shows a single fixed page explaining what's going on
func startIntroServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(introPage))
	})
	panic(http.ListenAndServe(":" + settings["port"], nil))
}

// fileExists returns true if the specified file or directory exists
func fileExists(name string) bool {
	if name != "" {
		if fd, err := os.Open(name); err == nil {
			fd.Close()
			return true
		}
	}
	return false
}

// loadSettings parses a settings file, if it exists, to configure some basic
// application settings, such as where the app/ and data/ directories are.
// For "fooBarDir", it also allows overriding via a "FOO_BAR_DIR" env variable.
func loadSettings(filename string) {
	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		if line != "" && !strings.HasPrefix(line, "#") {
			fields := strings.SplitN(line, "=", 2)
			if len(fields) != 2 {
				log.Fatalln("cannot parse settings:", scanner.Text())
			}
			key := strings.Trim(fields[0], " \t")
			value := strings.Trim(fields[1], " \t")
			env := os.Getenv(key)
			if env != "" {
				value = env
			}
			settings[capsToAllCaps(key)] = value
		}
	}
}

// capsToAllCaps converts "fooBarDir" to "FOO_BAR_DIR"
func capsToAllCaps(s string) (result string) {
	s = strings.ToLower(s)
	t := strings.Split(s, "_")
	for i, _ := range t {
		result += strings.ToUpper(t[i][:1]) + t[i][1:]
	}
	return
}

// introPage contains the HTML shown when the application cannot start normally
const introPage = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Welcome to HouseMon</title>
  </head>
  <body>
    <blockquote>
      <h3>Welcome to HouseMon</h3>
      <p>Whoops ... the main application files were not found.</p>
      <p>Please launch this application from the HouseMon directory.</p>
    </blockquote>
    <script>
      setInterval(function() {
        ws = new WebSocket("ws://" + location.host + "/ws");
        ws.onopen = function() {
          window.location.reload(true)
        }
      }, 1000)
    </script>
  </body>
</html>`
