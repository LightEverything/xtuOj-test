package main

import "xtuOj/router"

func main() {
	r := router.Router()
	r.Run(":80")
}
