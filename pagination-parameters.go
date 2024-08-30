package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PageQueryParams struct {
	startAt int64
	count   int64
}

func pageQueryFromRequestQueryParams(r *http.Request) (PageQueryParams, error) {
	params := mux.Vars(r)

	startAt := params["startAt"]
	count := params["count"]

	if startAt == "" {
		startAt = "0"
	}

	if count == "" {
		count = "250"
	}

	return paramsToQueryParams(startAt, count)
}

func paramsToQueryParams(startAt string, count string) (PageQueryParams, error) {
	var params PageQueryParams

	startAtInt, err := strconv.Atoi(startAt)
	if err != nil {
		return params, err
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		return params, err
	}

	if countInt < 0 || startAtInt < 0 {
		return params, fmt.Errorf("pagination parameters cannot be negative")
	}

	params.count = int64(countInt)
	params.startAt = int64(startAtInt)

	return params, nil
}
