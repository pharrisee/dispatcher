package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
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
)

func main() {
	b, err := os.ReadFile("scripts.json")
	if err != nil {
		log.Fatal(err)
	}
	var inscripts inScripts
	err = json.Unmarshal(b, &inscripts)
	if err != nil {
		log.Fatal(err)
	}

	var scripts Scripts
	for _, s := range inscripts {
		i, err := time.ParseDuration(s.Interval)
		if err != nil {
			log.Fatal("converting interval: %w", err)
		}
		scripts = append(scripts, Script{Name: s.Name, Interval: i, NextRun: time.Now().Add(i)})
	}

	tick := time.NewTicker(time.Second)
	for t := range tick.C {
		for i, script := range scripts {
			if script.NextRun.Before(t) {
				script.NextRun = script.NextRun.Add(time.Duration(script.Interval))
				cmd := exec.Command("/bin/bash", "-c", "./"+script.Name)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					log.Printf("%s: %s", script.Name, err)
				}
			}
			scripts[i] = script
		}
	}
	select {}
}
