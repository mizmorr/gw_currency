package main

import "github.com/mizmorr/gw_currency/gw-exchanger/internal/config"

func main() {
	cfg := config.Get()
	cfg.Print()
}
