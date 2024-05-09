/*
	Copyright (C) 2024  Pagefault Games

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/gob"
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/noahwillcrow/myfork-rogueserver/api"
	"github.com/noahwillcrow/myfork-rogueserver/db"
)

func main() {
	// flag stuff
	debug := flag.Bool("debug", false, "use debug mode")

	proto := flag.String("proto", "tcp", "protocol for api to use (tcp, unix)")
	addr := flag.String("addr", "0.0.0.0", "network address for api to listen on")
	passkey := flag.String("passkey", "passkey", "passkey for requests to authorize access to the api")

	dbuser := flag.String("dbuser", "pokerogue", "database username")
	dbpass := flag.String("dbpass", "", "database password")
	dbproto := flag.String("dbproto", "tcp", "protocol for database connection")
	dbaddr := flag.String("dbaddr", "localhost", "database address")
	dbname := flag.String("dbname", "pokeroguedb", "database name")

	allowedOrigins := flag.String("allowedorigins", "*", "CORS allowed origins")

	statRefreshCronSpec := flag.String("statrefreshcronspec", "@hourly", "CRON spec for stat refresh job")

	flag.Parse()

	// register gob types
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})

	// get database connection
	err := db.Init(*dbuser, *dbpass, *dbproto, *dbaddr, *dbname)
	if err != nil {
		log.Fatalf("failed to initialize database: %s", err)
	}

	err = db.PrepareTables()
	if err != nil {
		log.Fatalf("failed to prepare tables: %s", err)
	}

	// create listener
	listener, err := createListener(*proto, *addr)
	if err != nil {
		log.Fatalf("failed to create net listener: %s", err)
	}

	mux := http.NewServeMux()

	// init api
	api.Init(mux, *statRefreshCronSpec)

	log.Printf("rogueserver ready to start on %s://%s", *proto, *addr)

	// start web server
	if *debug {
		err = http.Serve(listener, debugHandler(mux))
	} else {
		err = http.Serve(listener, standardRequestHandler(mux, *allowedOrigins, *passkey))
	}
	if err != nil {
		log.Fatalf("failed to create http server or server errored: %s", err)
	}
}

func createListener(proto, addr string) (net.Listener, error) {
	if proto == "unix" {
		os.Remove(addr)
	}

	listener, err := net.Listen(proto, addr)
	if err != nil {
		return nil, err
	}

	if proto == "unix" {
		os.Chmod(addr, 0777)
	}

	return listener, nil
}

func debugHandler(router *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		router.ServeHTTP(w, r)
	})
}

func standardRequestHandler(router *http.ServeMux, allowedOrigins string, passkey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("x-passkey") != passkey {
			w.WriteHeader(http.StatusForbidden)
		}

		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		router.ServeHTTP(w, r)
	})
}
