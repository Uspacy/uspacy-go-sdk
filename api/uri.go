package api

import "os"

var mainHost = os.Getenv("HOST")

var emptyHeaders map[string]string
