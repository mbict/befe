package main

type entry struct {
	JS  string
	CSS string
}

var distEntryPoints = map[string]entry{
	"account":    {JS: "index-f92b2378.js", CSS: "index-94d9dd7e.css"},
	"backoffice": {JS: "index-2109d681.js", CSS: "index-20120fef.css"},
	"agenda":     {JS: "index-e2a372f0.js", CSS: "index-30593df1.css"},
	"admin":      {JS: "index-45ed39f0.js", CSS: "index-392872ec.css"},
}
