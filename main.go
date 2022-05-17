package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type (
	inScript struct {
		Name     string `json:"name"`
		Interval string `json:"interval"`
	}
	Script struct {
		Name     string        `json:"name"`
		Interval time.Duration `json:"interval"`
		NextRun  time.Time     `json:"-"`
	}
	Scripts   []Script
	inScripts []inScript
	inDisp    struct {
		Scripts  inScripts `json:"scripts"`
		Interval string    `json:"overall-interval"`
	}
	Disp struct {
		Scripts  Scripts `json:"scripts"`
		Interval time.Duration
	}
)

func main() {
	id := loadInDisp()
	disp := convertInDisp(id)
	runScripts(disp)
}

func convertInDisp(id inDisp) Disp {
	disp := Disp{}

	for _, s := range id.Scripts {
		i, err := time.ParseDuration(s.Interval)
		if err != nil {
			log.Fatal("converting interval: %w", err)
		}
		disp.Scripts = append(disp.Scripts, Script{Name: s.Name, Interval: i, NextRun: time.Now().Add(i)})
	}

	var err error
	disp.Interval, err = time.ParseDuration(id.Interval)
	if err != nil {
		log.Fatal("converting overall interval: %w", err)
	}
	return disp
}

func loadInDisp() inDisp {
	b, err := os.ReadFile("dispatcher.json")
	if err != nil {
		log.Fatal(err)
	}
	id := inDisp{}
	err = json.Unmarshal(b, &id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func runScripts(disp Disp) {
	tick := time.NewTicker(disp.Interval)
	for t := range tick.C {
		for i, script := range disp.Scripts {
			go func(i int, script Script) {
				if script.NextRun.Before(t) {
					now := time.Now()
					script.NextRun = script.NextRun.Add(time.Duration(script.Interval))
					cmd := exec.Command("/bin/bash", "-c", "./"+script.Name)
					out, err := cmd.CombinedOutput()
					if err != nil {
						log.Printf("%s: %s", script.Name, err)
						return
					}
					outS := strings.TrimSpace(string(out))
					dur := time.Since(now)
					fmt.Printf("%s :: %s (%s): %s\n", t.Format("2006-01-02 15:04:05.000"), script.Name, dur, outS)
				}
				disp.Scripts[i] = script
			}(i, script)
		}
	}
	select {}

}
