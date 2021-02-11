package command

import (
	"github.com/tendermint/tendermint/libs/json"
	"testing"
)

func TestCommand(t *testing.T) {
	data := `{"recipient":"0xDB2Ae41912f8c53fd6e5475a1432dA65B4A09127","type":"send_to_eth","fee":"1"}`

	cmd := Command{}
	if err := json.Unmarshal([]byte(data), &cmd); err != nil {
		t.Fatalf("Unmarshalling failed: %s", err.Error())
	}

	if err := cmd.Validate(); err != nil {
		t.Fatalf("Validation failed: %s", err.Error())
	}
}
