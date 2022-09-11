package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func jsonToMap(jsonStr string) map[string]interface{} {
	result := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &result)
	return result
}

func fetchData(url string) <-chan map[string]interface{} {

	r := make(chan map[string]interface{})

	go func() {
		defer close(r)

		response, err := http.Get(url)

		if err != nil {
			fmt.Printf(err.Error())
			r <- make(map[string]interface{})
			return
		}

		var parsed map[string]interface{}
		body, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Printf(err.Error())
			r <- make(map[string]interface{})
			return
		}

		err = json.Unmarshal([]byte(body), &parsed)

		r <- parsed
	}()

	return r
}
