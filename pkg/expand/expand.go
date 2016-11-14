package expand

import (
	"encoding/base64"
	"fmt"
	"github.com/golang/glog"
	"regexp"
)

func DoExpand(src []byte, values map[string]interface{}) ([]byte, error) {
	var err error

	// All
	expr := `\$(\({1,2})([[:alnum:]_\.\-]+)(\|base64)?\){1,2}|(\{{2})([[:alnum:]_\.\-]+)(\|base64)?\}{2}`
	re := regexp.MustCompile(expr)
	expandFunction := func(match []byte) []byte {
		re := regexp.MustCompile(expr)

		matchStr := string(match[:])
		result := re.FindStringSubmatch(matchStr)

		if result[0] != matchStr {
			glog.Fatalf("Unexpected match: %q", matchStr)
		}

		if result[2] == "" && result[5] == "" {
			glog.Fatalf("No variable defined within: %q", matchStr)
		}

		key := result[2] + result[5]
		replacement := values[key]

		if replacement == nil {
			err = fmt.Errorf("Key not found: %q", key)
			return match
		}

		if (result[3] + result[6]) == "|base64" {
			replacement = base64.StdEncoding.EncodeToString([]byte(replacement.(string)))
		}

		var s string
		delim := result[1] + result[4]
		switch len(delim) {
		case 1:
			s = fmt.Sprintf("\"%v\"", replacement)
		case 2:
			s = fmt.Sprintf("%v", replacement)
		default:
			glog.Fatalf("Unexpected delimiter %q count: %q", delim, len(delim))
		}

		return []byte(s)
	}

	expanded := re.ReplaceAllFunc(src, expandFunction)
	if err != nil {
		return nil, err
	}

	return expanded, nil
}