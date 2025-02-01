package utils

import (
	"bufio"
	"os"
	"path/filepath"
)

func GetScripts() (temps string, individuals string, benefits string) {
	currentDir, _ := os.Getwd()
	dirScripts := filepath.Join(currentDir, "scripts")
	scripts, err := os.ReadDir(dirScripts)
	if err != nil {
		panic(err)
	}
	var (
		t, i, b string
	)
	for _, script := range scripts {
		switch script.Name() {
		case "temps.sql":
			t = filepath.Join(dirScripts, script.Name())
		case "individuals.sql":
			i = filepath.Join(dirScripts, script.Name())
		case "benefits.sql":
			b = filepath.Join(dirScripts, script.Name())
		}
	}
	temps = ScanScript(t)
	individuals = ScanScript(i)
	benefits = ScanScript(b)
	return temps, individuals, benefits
}

func ScanScript(file string) string {
	var script string
	str, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(str)
	for scanner.Scan() {
		script += scanner.Text()
	}
	return script
}
