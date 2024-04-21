package encoding

import "encoding/json"

func ToString(input any) string {
	data, err := json.Marshal(input)
	if err != nil {
		return ""
	}

	return string(data)
}
