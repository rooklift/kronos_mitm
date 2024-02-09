package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s [nodes] [path] [...additional args]", os.Args[0])
		return
	}

	nodes_limit, err := strconv.Atoi(os.Args[1])
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
	adjuster_mitm(os.Stdin, i_pipe, nodes_limit)
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
		if strings.HasPrefix(s, "go ") {
			output.Write([]byte(fmt.Sprintf("go nodes %d\n", nodes_limit)))
		} else {
			output.Write([]byte(s))
			output.Write([]byte{'\n'})
		}
	}
}
