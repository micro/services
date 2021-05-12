package handler

import (
	"testing"
)

func TestBase64Decode(t *testing.T) {
	_, _, err := base64ToImage("data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==")
	if err != nil {
		t.Fatal(err)
	}
}
