package app

import (
	"log"
	"net/http"
	"net/http/pprof"
	"time"
)

func (app *Payssage) startServer(mux http.Handler) {
	srv := &http.Server{
		Addr:         ":" + app.cfg.Port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Printf("%s listening on :%s", app.cfg.AppName, app.cfg.Port)
	log.Fatal(srv.ListenAndServe())
}

func servePprof(port string) {
	pmux := http.NewServeMux()
	pmux.HandleFunc("/debug/pprof/", pprof.Index)
	pmux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	pmux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	pmux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	pmux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      pmux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Printf("payssage pprof listening on :%s", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("payssage pprof server error: %v", err)
	}
}
