package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	cfgFile = flag.String("config", "./StaticSite", "The path to your sites configuration file")
)

var cfg *Config

func renderPage(page string, w io.Writer) {
	templates, err := filepath.Glob(filepath.Join(cfg.TemplateDir, "*.tmpl"))
	if err != nil {
		log.Fatal(err)
	}

	files := append(templates, page)
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(fmt.Errorf("Couldn't parse templates '%s': %v", strings.Join(files, ","), err))
	}
	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	var err error
	cfg, err = loadConfig(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Removing destination directory '%s'", cfg.OutputDir)
	if err := os.RemoveAll(cfg.OutputDir); err != nil {
		log.Fatal(err)
	}

	log.Printf("Copying static files to output dir")
	if err := copyDir(cfg.StaticDir, cfg.OutputDir); err != nil {
		log.Fatal(err)
	}

	pages, err := filepath.Glob(filepath.Join(cfg.PageDir, "*"))
	if err != nil {
		log.Fatal(err)
	}
	for _, page := range pages {
		outfile, err := filepath.Rel(cfg.PageDir, page)
		if err != nil {
			log.Panic(err)
		}
		outfile = filepath.Join(cfg.OutputDir, outfile)
		f, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Couldn't open output file '%s': %v", outfile, err)
		}

		log.Printf("Rendering page: %s to %s...\n", page, outfile)
		renderPage(page, f)
	}

}

func copyFile(source string, dest string) error {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()

	si, err := sf.Stat()
	if err != nil {
		return err
	}

	df, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, si.Mode())
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	if err != nil {
		return err
	}

	return nil
}

func copyDir(source string, dest string) error {
	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("Source is not a directory")
	}

	// ensure dest dir does not already exist
	_, err = os.Open(dest)
	if !os.IsNotExist(err) {
		return fmt.Errorf("Destination already exists")
	}

	// create dest dir
	if err := os.MkdirAll(dest, fi.Mode()); err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(source)
	for _, entry := range entries {
		sfp := source + "/" + entry.Name()
		dfp := dest + "/" + entry.Name()
		if entry.IsDir() {
			err = copyDir(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		} else {
			// perform copy
			err = copyFile(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func noExt(path string) string {
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[:i]
		}
	}
	return ""
}
