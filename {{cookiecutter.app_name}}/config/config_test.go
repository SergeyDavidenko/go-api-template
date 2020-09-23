package config

import (
	"testing"
)

func TestLoadDefaultConfig(t *testing.T) {
	config, _ := LoadConf("")
	if config.Core.Lable != "{{cookiecutter.app_name}}" {
		t.Errorf("lable not eq {{cookiecutter.app_name}}, got %s", config.Core.Lable)
	}
	if config.Log.Level != "info" {
		t.Errorf("log level not eq info, got %s", config.Log.Level)
	}
}

func TestLoadConfig(t *testing.T) {
	_, err := LoadConf("/asdasd/asdasd")
	if err == nil {
		t.Error("Not return err on not exist file")
	}

}