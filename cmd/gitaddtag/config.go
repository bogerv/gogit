package main

type Config struct {
	Paths    []string `json:"paths"`    // project absolute path (e.g. D:\golang\otc)
	Branches []string `json:"branches"` // project branches
	Version  string   `json:"version"`  // tag version
	Message  string   `json:"message"`  // tag message
}
