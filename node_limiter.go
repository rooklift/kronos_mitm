package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	Nodes int
	Path string
	Args []string
}

func main() {

	mypath, err := os.Executable()
	if err != nil {
		fmt.Printf("Couldn't find myself\n")
		return
	}

	file, err := ioutil.ReadFile(filepath.Join(filepath.Dir(mypath), "config.json"))
	if err != nil {
		fmt.Printf("Failed to load config file: %v\n", err)
		return
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Printf("Couldn't parse config file: %v\n", err)
		return
	}

	exec_command := exec.Command(config.Path, config.Args...)
	i_pipe, _ := exec_command.StdinPipe()
	o_pipe, _ := exec_command.StdoutPipe()
	e_pipe, _ := exec_command.StderrPipe()

	err = exec_command.Start()
	if err != nil {
		fmt.Printf("Failed to start engine: %v\n", err)
		return
	}

	go mitm(o_pipe, os.Stdout)
	go mitm(e_pipe, os.Stderr)
	adjuster_mitm(os.Stdin, i_pipe, config.Nodes)
}

func mitm(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		output.Write(scanner.Bytes())
		output.Write([]byte{'\n'})
	}
}

func adjuster_mitm(input io.Reader, output io.Writer, nodes_limit int) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.HasPrefix(s, "go") {
			output.Write([]byte(fmt.Sprintf("go nodes %d\n", nodes_limit)))
		} else {
			output.Write([]byte(s))
			output.Write([]byte{'\n'})
		}
	}
}
