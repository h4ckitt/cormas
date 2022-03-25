package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

func GetTagUID(tag interface{}) (string, error) {
	var (
		q         string
		variables = make(map[string]string)
	)

	fmt.Println("Beginning Of Func")

	if x, ok := tag.(map[string]interface{}); ok {
		fmt.Println("Inside Models Func")
		if len(x) == 0 || x["name"].(string) == "" {
			return "", errors.New("invalid tag")
		}

		q =
			`
		query Tags($name: string) {
			tag(func: type(HashTag)) @filter(eq(name, $name)) {
				uid
			}
		}			
		`

		variables["$name"] = x["name"].(string)
	} /*else if x, ok := tag.(string); ok {		//I'll Activate This If I Find The Need, For Now, It's Cool
		fmt.Println("Inside UID Func")
		if x == "" {
			return "", errors.New("invalid tag")
		}

		q =
			`
		query Tags($uid: string) {
			tag(func: uid($uid)) {
				uid
			}
		}
		`

		variables["$uid"] = x

	}*/
	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, variables)

	if err != nil {
		return "", err
	}

	result := struct {
		Result []struct {
			UID string `json:"uid"`
		} `json:"tag"`
	}{}

	json.Unmarshal(resp.Json, &result)

	if len(result.Result) == 0 {
		/*msg := "tag doesn't exist"
		if _, ok := tag.(map[string]interface{}); ok {
			msg = "tag doesn't exist yet"
		}*/
		return "", errors.New("tag doesn't exist yet")
	}

	return result.Result[0].UID, nil
}
