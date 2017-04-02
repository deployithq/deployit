//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package editor

import (
	"bufio"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"os/exec"
)

const default_editor = "vim"

type Editor struct {
	context   Context
	SplitFunc bufio.SplitFunc
}

func Run(content io.Reader) (*Editor, error) {

	var editor = &Editor{SplitFunc: bufio.ScanLines}

	err := editor.newContext(content)
	if err != nil {
		return nil, err
	}

	return editor, nil
}

func (e *Editor) newContext(content io.Reader) error {

	var err error

	err = e.context.Write(content, e.SplitFunc)
	if err != nil {
		return err
	}

	cmd := e.cmd(e.context.Name)
	err = cmd.Run()
	if err != nil {
		return err
	}

	err = e.context.Read(e.context.Name, e.SplitFunc)
	if err != nil {
		return err
	}

	return nil
}

func (Editor) cmd(filename string) *exec.Cmd {

	var path = os.Getenv("EDITOR")

	if path == "" {
		path = default_editor
	}

	cmd := exec.Command(path, filename)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func (e *Editor) Lines() []string {
	return e.context.output.c
}

func (e *Editor) Bytes() []byte {
	return e.context.output.Bytes()
}

func (e *Editor) Line(i int) string {
	return e.context.output.c[i]
}

func (e *Editor) FromYAML(i interface{}) error {
	return yaml.Unmarshal(e.context.output.Bytes(), i)
}

func (e *Editor) FromJSON(i interface{}) error {
	return json.Unmarshal(e.context.output.Bytes(), i)
}
