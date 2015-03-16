package loghttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func outputBody(original io.ReadCloser, headers http.Header) (io.ReadCloser, error) {
	defer original.Close()

	var bs bytes.Buffer
	_, err := io.Copy(&bs, original)
	if err != nil {
		return nil, err
	}

	contentType := headers.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		formatJSON(bs.Bytes())
	} else {
		fmt.Println(bs.String())
	}

	return ioutil.NopCloser(strings.NewReader(bs.String())), nil
}

func formatJSON(raw []byte) {
	var data map[string]interface{}

	err := json.Unmarshal(raw, &data)
	if err != nil {
		fmt.Printf("Unable to parse JSON: %s\n\n", err)
		fmt.Println(string(raw))

		return
	}

	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Unable to re-marshal JSON: %s\n", err)
		fmt.Println(string(raw))

		return
	}

	fmt.Println(string(pretty))
}
