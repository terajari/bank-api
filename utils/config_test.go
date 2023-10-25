package utils

import "testing"

func TestLoadConfig(t *testing.T) {
	driver := "postgres"
	source := "postgresql://postgres:1234@localhost:5432/bank_api?sslmode=disable"

	got, err := LoadConfig("../.env")
	if err != nil {
		t.Errorf("LoadConfig() error = %v", err)
		return
	}
	if got.DBDriver != driver {
		t.Errorf("got.DBDriver = %v, want %v", got.DBDriver, driver)
	}
	if got.DBSource != source {
		t.Errorf("got.DBSource = %v, want %v", got.DBSource, source)
	}
}
