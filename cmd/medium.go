package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	medium "github.com/medium/medium-sdk-go"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

func getMediumClient() (*medium.Medium, error) {
	config, err := toml.LoadFile("medium.toml")
	if err != nil {
		return nil, errors.Wrap(err, "failed to load medium.toml")
	}
	token := config.Get("user.token")
	if token == nil {
		return nil, errors.New("user.token is not in medium.toml.")
	}
	if token == "" {
		return nil, errors.New("token is empty")
	}
	m := medium.NewClientWithAccessToken(token.(string))
	return m, nil
}

func getWorkspace() (string, error) {
	config, err := toml.LoadFile("medium.toml")
	if err != nil {
		return "", errors.Wrap(err, "failed to load medium.toml")
	}

	dir := config.Get("user.workspace")
	if dir == nil {
		return "", errors.New("user.workspace is not in medium.toml")
	}
	return dir.(string), nil
}

func getAuthorID() (string, error) {
	config, err := toml.LoadFile("medium.toml")
	if err != nil {
		return "", errors.Wrap(err, "failed to load medium.toml")
	}

	if !config.Has("user.author_id") {
		return "", nil
	}

	return config.Get("user.author_id").(string), nil
}

func getUserName() (string, error) {
	config, err := toml.LoadFile("medium.toml")
	if err != nil {
		return "", errors.Wrap(err, "failed to load medium.toml")
	}

	if !config.Has("user.name") {
		return "", nil
	}

	return config.Get("user.name").(string), nil
}

func setAuthorID(s string) error {
	config, err := toml.LoadFile("medium.toml")
	if err != nil {
		return errors.Wrap(err, "failed to load medium.toml")
	}

	config.Set("user.author_id", s)
	f, err := os.Open("medium.toml")
	if err != nil {
		return errors.Wrap(err, "failed to open medium.toml")
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	config.WriteTo(f)
	return nil
}

func getPostOptions(t *toml.TomlTree, p string) (*medium.CreatePostOptions, error) {
	requiredFields := []string{"title", "tags", "format", "content_path"}
	for _, v := range requiredFields {
		if !t.Has("info." + v) {
			return nil, errors.Errorf("toml doesn't have field %s", v)
		}
		fmt.Println("[Pass] info." + v)
	}

	b, err := ioutil.ReadFile(filepath.Join(p, t.Get("info.content_path").(string)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read content file")
	}

	rawTags := t.Get("info.tags").([]interface{})
	tags := make([]string, len(rawTags))
	for i, v := range rawTags {
		tags[i] = v.(string)
	}

	po := &medium.CreatePostOptions{
		Title:         t.Get("info.title").(string),
		Tags:          tags,
		ContentFormat: medium.ContentFormat(t.Get("info.format").(string)),
		Content:       string(b),
	}

	if t.Has("info.CanonicalURL") {
		po.CanonicalURL = t.Get("info.CanonicalURL").(string)
	}

	if t.Has("info.PublishStatus") {
		po.PublishStatus = medium.PublishStatus(t.Get("info.PublishStatus").(string))
	}

	if t.Has("info.License") {
		po.License = medium.License(t.Get("info.License").(string))
	}
	return po, nil
}
