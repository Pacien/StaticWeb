/*

	This file is part of StaticWeb <https://github.com/Pacien/StaticWeb>.

	StaticWeb is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	StaticWeb is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with StaticWeb. If not, see <http://www.gnu.org/licenses/>.

*/

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var params struct {
	addr, port, dir, log string
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	host := strings.Split(r.Host, ":")

	if host[0] == "" {
		log.Println("Undefined host")
		http.NotFound(w, r)
		return
	}

	request := r.URL.Path[1:]
	requestedFile := params.dir + "/" + host[0] + "/" + request
	log.Println(requestedFile)

	file, err := os.Stat(requestedFile)
	if err != nil {
		log.Println(err)
		if os.IsNotExist(err) {
			http.NotFound(w, r)
		} else if os.IsPermission(err) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if file.IsDir() && !strings.HasSuffix(requestedFile, "/") {
		http.Redirect(w, r, r.URL.Path+"/", http.StatusFound)
		return
	}

	http.ServeFile(w, r, requestedFile)
}

func init() {
	flag.StringVar(&params.addr, "addr", "127.0.0.1", "Address to listen.")
	flag.StringVar(&params.port, "port", "8080", "Port to listen.")
	flag.StringVar(&params.dir, "dir", ".", "Absolute or relative path to the root directory to serve.")
	flag.StringVar(&params.log, "log", "", "Absolute or relative path to the log file. Leave empty for stdout.")
	flag.Parse()
}

func main() {
	fmt.Println("StaticWeb <https://github.com/Pacien/StaticWeb>")

	if params.log != "" {
		logFile, err := os.OpenFile(params.log, os.O_WRONLY, 0666)
		if os.IsNotExist(err) {
			log.Println("Log file not found, creating a new log file:", err)
			logFile, err = os.Create(params.log)
			if err != nil {
				log.Println("Cannot create log file:", err)
				return
			}
		} else if err != nil {
			log.Println("Cannot open log file:", err)
			return
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	log.Println("Listening on " + params.addr + ":" + params.port)

	http.HandleFunc("/", defaultHandler)
	err := http.ListenAndServe(params.addr+":"+params.port, nil)
	if err != nil {
		log.Println(err)
		return
	}
}
