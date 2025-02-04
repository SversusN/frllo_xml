package utils

import (
	"bufio"
	"os"
	"path/filepath"
)

func GetScripts() (temps string, individuals string, benefits string, recipes string) {
	currentDir, _ := os.Getwd()
	dirScripts := filepath.Join(currentDir, "scripts")
	scripts, err := os.ReadDir(dirScripts)
	if err != nil {
		panic(err)
	}
	var (
		t, i, b, r string
	)
	for _, script := range scripts {
		switch script.Name() {
		case "temps.sql":
			t = filepath.Join(dirScripts, script.Name())
		case "individuals.sql":
			i = filepath.Join(dirScripts, script.Name())
		case "benefits.sql":
			b = filepath.Join(dirScripts, script.Name())

		case "recipes.sql":
			r = filepath.Join(dirScripts, script.Name())
		}
	}
	temps = ScanScript(t)
	individuals = ScanScript(i)
	benefits = ScanScript(b)
	recipes = ScanScript(r)
	return temps, individuals, benefits, recipes
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
