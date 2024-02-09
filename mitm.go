package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
)

var nodes_limit int

func main() {

	var err error

	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s [nodes] [path] [...additional args]", os.Args[0])
		return
	}

	nodes_limit, err = strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid nodes arg: %v", err)
		return
	}

	exec_command := exec.Command(os.Args[2], os.Args[3:]...)
	i_pipe, _ := exec_command.StdinPipe()
	o_pipe, _ := exec_command.StdoutPipe()
	e_pipe, _ := exec_command.StderrPipe()

	err = exec_command.Start()
	if err != nil {
		fmt.Printf("Failed to start engine: %v", err)
		return
	}

	go mitm(o_pipe, os.Stdout)
	go mitm(e_pipe, os.Stderr)
	adjuster_mitm(os.Stdin, i_pipe)
}

func mitm(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		output.Write(scanner.Bytes())
		output.Write([]byte{'\n'})
	}
}

func adjuster_mitm(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		output.Write(scanner.Bytes())
		output.Write([]byte{'\n'})
	}
}
