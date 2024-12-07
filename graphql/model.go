package main

type Account struct {
	ID     string  `json:"id"`
	Name   string  `json:"naem"`
	Orders []Order `json:"orders"`
}
