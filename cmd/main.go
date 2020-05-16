package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/popsu/mtgpwscraper/pwscraper"
)

func main() {
	conf, output, err := parseFlags(os.Args[0], os.Args[1:])
	if err == flag.ErrHelp {
		fmt.Println(output)
		os.Exit(2)
	} else if err != nil {
		fmt.Println("got error:", err)
		fmt.Println("output:\n", output)
		os.Exit(1)
	}

	pwscraper.Execute(conf)
}

func parseFlags(progname string, args []string) (config *pwscraper.Config, output string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	conf := pwscraper.NewDefaultConfig()
	flags.StringVar(&conf.DciNumber, "dcinumber", "", "DCI Number")
	flags.StringVar(&conf.PwpCookieValue, "cookie", "", "Cookie named PWP.ASPXAUTH in wizards.com site while logged in")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}
	return conf, buf.String(), nil
}
