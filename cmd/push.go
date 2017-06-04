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

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "publish your article to medium",
	Long: `
`,
	RunE: publishArticle,
}

func init() {
	RootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func publishArticle(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("args is not only one.")
	}

	// info.toml取ってくる
	ws, err := getWorkspace()
	if err != nil {
		return errors.Wrap(err, "failed to get workspace path")
	}

	article := filepath.Join(ws, args[0])
	_, err = os.Stat(article)
	if err != nil {
		return errors.Wrap(err, "failed to get dir to post")
	}

	config, err := toml.LoadFile(filepath.Join(article, "info.toml"))
	if err != nil {
		return errors.Wrap(err, "failed to load infomation toml")
	}

	// 必要な情報を取得
	// info.toml から情報を引っ張ってきてPostできるように整形
	cpo, err := getPostOptions(config, article)
	if err != nil {
		return errors.Wrap(err, "failed to get information from toml object")
	}

	// UserID (見た目に出るIDとはべつのもの) を取得する
	m, err := getMediumClient()
	if err != nil {
		return errors.Wrap(err, "failed to create Medium Client")
	}
	u, err := m.GetUser("")
	if err != nil {
		return errors.Wrap(err, "failed to get user infomation from medium")
	}
	cpo.UserID = u.ID
	fmt.Printf("[Post] you post new Story as %s", u.Username)

	// 実際にポストする
	p, err := m.CreatePost(*cpo)
	if err != nil {
		return errors.Wrap(err, "failed to post content to medium")
	}

	// 成功したらURLを出力して終了
	fmt.Printf("[Post] suceed to post\nURL: %s\n", p.URL)
	return nil
}
