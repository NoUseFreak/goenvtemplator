package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Masterminds/sprig"
	"io/ioutil"
	"log"
	"path/filepath"
	"text/template"
)

func Require(arg string) (string, error) {
	if len(arg) == 0 {
		return "", errors.New("Required argument is missing!")
	}
	return arg, nil
}

var funcMap = template.FuncMap{
	"require": Require,
}

func generateTemplate(source, name string, delimLeft string, delimRight string) (string, error) {
	var t *template.Template
	var err error
	t, err = template.New(name).Delims(delimLeft, delimRight).Option("missingkey=error").Funcs(funcMap).Funcs(sprig.TxtFuncMap()).Parse(source)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	// hacking because go 1.7 fails to throw error, see https://github.com/golang/go/commit/277bcbbdcd26f2d64493e596238e34b47782f98e
	emptyHash := map[string]interface{}{}
	if err = t.Execute(&buffer, &emptyHash); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func generateFile(templatePath, destinationPath string, debugTemplates bool, delimLeft string, delimRight string) error {
	if !filepath.IsAbs(templatePath) {
		return fmt.Errorf("Template path '%s' is not absolute!", templatePath)
	}

	if !filepath.IsAbs(destinationPath) {
		return fmt.Errorf("Destination path '%s' is not absolute!", destinationPath)
	}

	var slice []byte
	var err error
	if slice, err = ioutil.ReadFile(templatePath); err != nil {
		return err
	}
	s := string(slice)
	result, err := generateTemplate(s, filepath.Base(templatePath), delimLeft, delimRight)
	if err != nil {
		return err
	}

	if debugTemplates {
		log.Printf("Printing parsed template to stdout. (It's delimited by 2 character sequence of '\\x00\\n'.)\n%s\x00\n", result)
	}

	if err = ioutil.WriteFile(destinationPath, []byte(result), 0664); err != nil {
		return err
	}

	return nil
}
