package cli

import (
	"encoding/json"
	"fmt"
	"github.com/drewsonne/go-gocd/gocd"
	"os"
)

const GoCDUtilityDefaultVersion = "dev"

func Version() string {
	if tag := os.Getenv("TAG"); tag != "" {
		return tag
	} else {
		if commit := os.Getenv("COMMIT"); commit != "" {
			return commit[0:8]
		} else {
			return GoCDUtilityDefaultVersion
		}
	}
}

func cliAgent() *gocd.Client {
	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}

	return cfg.Client()
}

func handeErrOutput(reqType string, err error) error {
	return handleOutput(nil, nil, reqType, err)
}

func handleOutput(r interface{}, hr *gocd.APIResponse, reqType string, err error) error {
	var b []byte
	var o map[string]interface{}
	if err != nil {
		o = map[string]interface{}{
			"Error": err.Error(),
		}
	} else if hr.HTTP.StatusCode >= 200 && hr.HTTP.StatusCode < 300 {
		o = map[string]interface{}{
			fmt.Sprintf("%sResponse", reqType): r,
		}
		//} else if hr.HTTP.StatusCode == 404 {
		//	o = map[string]interface{}{
		//		"Error": fmt.Sprintf("Could not find resource for '%s' action.", reqType),
		//	}
	} else {

		b1, _ := json.Marshal(hr.HTTP.Header)
		b2, _ := json.Marshal(hr.Request.HTTP.Header)
		o = map[string]interface{}{
			"Error":           "An error occured while retrieving the resource.",
			"Status":          hr.HTTP.StatusCode,
			"ResponseHeader":  string(b1),
			"ResponseBody":    hr.Body,
			"RequestBody":     hr.Request.Body,
			"RequestEndpoint": hr.Request.HTTP.URL.String(),
			"RequestHeader":   string(b2),
		}
	}
	b, err = json.MarshalIndent(o, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

	return nil
}
