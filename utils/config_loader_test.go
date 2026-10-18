package utils

import (
	"fmt"
	"os"
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestConfigLoader(t *testing.T) {
	fmt.Println("========== TestConfigLoader ==========")
	data, err := os.ReadFile("../config/models.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}

	model := config.Models[config.DefaultModel]

	fmt.Println(model.Model)
	fmt.Println(model.BaseURL)
	fmt.Println(model.APIKeyEnv)
	fmt.Println(model.APIType)
}
