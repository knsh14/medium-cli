package cmd

import (
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
