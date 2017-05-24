// Copyright © 2017 Kenshi Kamata <kenshi.kamata@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"html/template"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const mediumTOML = `[[user]]
name = "{{.Name}}"
token = ""
workspace = "{{.Workspace}}"
`

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your medium workspace",
	Long: `
`,
	RunE: initialize,
}

func init() {
	RootCmd.AddCommand(initCmd)
}

func initialize(cmd *cobra.Command, args []string) error {
	fmt.Println("[init] Create toml file to set user info")
	err := createTOML()
	if err != nil {
		return errors.Wrap(err, "failed to create toml")
	}
	return nil
}

func createTOML() error {
	f, err := os.Create("medium.toml")
	if err != nil {
		return errors.Wrap(err, "failed to create medium.toml")
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	var name string
	fmt.Print("Username: ")
	fmt.Scanln(&name)

	currentDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to get current dir path")
	}

	UserInfo := struct {
		Name, Workspace string
	}{Name: name, Workspace: currentDir}

	// Create a new template and parse the letter into it.
	t := template.Must(template.New("mediumTOML").Parse(mediumTOML))
	if err := t.Execute(f, UserInfo); err != nil {
		return errors.Wrap(err, "failed to write template to toml")
	}

	return nil
}
