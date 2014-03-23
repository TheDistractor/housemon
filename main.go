package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/jcw/flow"
	jeebus "github.com/jcw/jeebus/gadgets"
)

var VERSION = "0.9.0" // can be adjusted by goxc

const defaults = `
# can be overridden with environment variables
APP_DIR = ./app
BASE_DIR = ./base
DATA_DIR = ./data
PORT = 5561
`

var config = jeebus.LoadConfig(defaults, "./config.txt")

func init() {
	flow.Registry["info"] = func() flow.Circuitry { return &infoCmd{} }
	jeebus.Help["info"] = `Show some basic information about HouseMon.`
}

type infoCmd struct{ flow.Gadget }

func (g *infoCmd) Run() {
	fmt.Printf("HouseMon %s + JeeBus %s + Flow %s\n\n",
		VERSION, jeebus.Version, flow.Version)
	flow.PrintRegistry()
	fmt.Println("\nUse 'help' for a list of commands or '-h' for a list of options.")
}

func main() {
	flag.Parse()

	if flag.NArg() > 0 {
		if factory, ok := flow.Registry[flag.Arg(0)]; ok {
			factory().Run()
			return
		}
		fmt.Fprintln(os.Stderr, "Unknown command:", flag.Arg(0), "(try 'help')")
		fmt.Fprintln(os.Stderr, "See http://jeelabs.net/projects/housemon/wiki")
		os.Exit(1)
	}

	fmt.Printf("Starting webserver for http://localhost:%s/\n", config["PORT"])

	// show intro page via a static webserver if the main app dir is absent
	fd, err := os.Open(config["APP_DIR"])
	if err != nil {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(introPage))
		})
		panic(http.ListenAndServe(":"+config["PORT"], nil))
	}
	fd.Close()

	// normal startup: save config info in database, then start the webserver
	setupDatabase()
	setupWebserver()
}

// database setup, save version and current config settings
func setupDatabase() {
	c := flow.NewCircuit()
	c.Add("db", "LevelDB")
	c.Add("sink", "Sink")
	c.Connect("db.Out", "sink.In", 0)
	c.Connect("db.Mods", "sink.In", 0)
	c.Feed("db.Name", config["DATA_DIR"])
	c.Feed("db.In", flow.Tag{"<clear>", "/config/"})
	for k, v := range config {
		c.Feed("db.In", flow.Tag{"/config/" + k, v})
	}
	c.Feed("db.In", flow.Tag{"/config/app", "HouseMon"})
	c.Feed("db.In", flow.Tag{"/config/version", VERSION})
	c.Run()
}

// webserver setup
func setupWebserver() {
	c := flow.NewCircuit()
	c.Add("http", "HTTPServer")
	c.Add("forever", "Forever") // run forever
	c.Feed("http.Handlers", flow.Tag{"/", config["APP_DIR"]})
	c.Feed("http.Handlers", flow.Tag{"/base", config["BASE_DIR"]})
	c.Feed("http.Handlers", flow.Tag{"/ws", "<websocket>"})
	c.Feed("http.Start", config["PORT"])
	c.Run()
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
