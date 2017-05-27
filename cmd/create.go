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
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const infoTemplate = `[[info]]
title = ""
tags = [""]
format = "markdown"
content_path = "content.md"
`

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create DIRNAME",
	Short: "create new directory to post",
	Long: `Create New Directory to write post.
DIRNAME
├── content.md
└── info.toml
`,
	RunE: createNewPostDir,
}

func init() {
	RootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func createNewPostDir(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.Errorf("args must be only one. got %d", len(args))
	}

	d, err := getWorkspace()
	if err != nil {
		return errors.Wrap(err, "failed to get workspace")
	}

	// create dir by args name
	article := filepath.Join(d, args[0])
	err = os.Mkdir(article, 0775)
	if err != nil {
		return errors.Wrapf(err, "failed to create directory %s", filepath.Join(d, args[0]))
	}

	// create toml and markdown in it
	info, err := os.Create(filepath.Join(article, "info.toml"))
	if err != nil {
		return errors.Wrap(err, "failed to create post infomation toml file")
	}
	defer func() {
		err := info.Close()
		if err != nil {
			panic(err)
		}
	}()
	_, err = fmt.Fprint(info, infoTemplate)
	if err != nil {
		return errors.Wrap(err, "failed to write template to info.toml")
	}

	c, err := os.Create(filepath.Join(article, "content.md"))
	if err != nil {
		return errors.Wrap(err, "failed to create content.md file")
	}
	defer func() {
		err := c.Close()
		if err != nil {
			panic(err)
		}
	}()

	fmt.Println("[create] finished")
	return nil
}
