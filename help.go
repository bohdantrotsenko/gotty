package main

var helpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} [options] <command> [<arguments...>]

VERSION:
   {{.Version}}{{if .Authors}}

AUTHOR:{{range .Authors}}
  {{.Name}}{{if .Email}} - <{{.Email}}>{{end}}{{end}}{{end}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`
