package model

import "net/http"

type Transaction struct {
	req *http.Request
	res http.ResponseWriter
}
