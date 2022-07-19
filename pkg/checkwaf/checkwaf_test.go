package checkwaf

import "testing"

func TestCheckWAF(t *testing.T) {
	CheckWAF("https://www.cloudflare.com/")
}
