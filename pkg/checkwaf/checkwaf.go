package checkwaf

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Reference: http://seclists.org/nmap-dev/2011/q2/att-1005/http-waf-detect.nse
// Reference: https://github.com/sqlmapproject/sqlmap/blob/c722f8e3bd4aac3a5dc2287db9e9dd04fb4ce257/lib/core/settings.py
const (
	IPSWAFCheckPayload = "AND 1=1 UNION ALL SELECT 1,NULL,'<script>alert(\"XSS\")</script>',table_name FROM information_schema.tables WHERE 2>1--/**/; EXEC xp_cmdshell('cat ../../../etc/passwd')#"
	letters            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func CheckWAF(rawURL string) (error, bool) {
	payload := fmt.Sprintf("%d %s", rand.Int(), IPSWAFCheckPayload)
	key := randStr(6)

	normalResp, err := http.DefaultClient.Get(rawURL)
	if err != nil {
		return err, false
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return err, false
	}
	query := u.Query()
	query.Add(key, payload)
	u.RawQuery = query.Encode()
	maliciousResp, err := http.DefaultClient.Get(u.String())
	if err != nil {
		return err, false
	}

	return nil, normalResp.StatusCode != maliciousResp.StatusCode
}
