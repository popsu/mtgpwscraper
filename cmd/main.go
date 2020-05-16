package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/popsu/mtgpwscraper/pwscraper"
)

const exampleCookie = "9A3D5E64320129F4764627A2E03F6025B4F838A2D40C5FE255809E2CC6F18EBCFAF0CA38A986460B7075183ACBF2385468C11F943BA4E4E2648B640F45C9B4E33643C150E664D52897F6DC35FC6DEBA0B3979E67497B3029FC48058388590A2C986FEC3F0894ACA1C90190064879ECEB8E836295803B4E25EB3C8A0EF0091090AAEF60A98463B8C68C0E7BBFE27958F9F2FE424474D752F1A413CCC3ABD8558CD3449EF92446C7733B23C1E9ABBB5FF240261F9D1579A910E074AF494A51C82D79399E158533F2697A2DC7286FF56B2636D6FAF67EBEFD983924DB1646487E2189F683DF2361A94B80EAE6A2B6C03A48FEC9ECCCD5DF200C7D2F08A606A3297B33B514942910D55D0D638D7C690FE18883331EBDC00F6F8CEF733720EFA46E5D5D1335205211EDDCF7952EDECC52183A5A66D3B7381149ABB32F467B985A547EB0E3A5DA3AFC26B0E0BDA548763C082FCA5523AC92DBCAF7A72D234E755B5F36989C714A87F9F4FAEBC0030384150780BC6758AB697637D4E7B318E6029F6FA183CCF8D62DD9B740F16248935C513484CBA91323EFCF61F7AE9E884B5FBA96A61FD87366A31C09473181F74637D06385B31BC390"
const loginURL = "https://www.wizards.com/Magic/PlaneswalkerPoints/"

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
	if conf.DciNumber == "" || conf.PwpCookieValue == "" {
		fmt.Printf("Error. Missing dcinumber and cookie arguments\n\n")
		fmt.Printf("Log in to %s to grab the cookie.\n", loginURL)
		fmt.Printf("In Chrome press Shift + CTRL +_J to open the developer console. Go to Application tab -> Storage -> Cookies -> https://www.wizards.com and in the list get the value of PWP.ASPXAUTH cookie to pass in as cookie to this program\n\n")
		fmt.Printf("Usage: %s -dcinumber <DCINUMBER> -cookie <COOKIE>\n\n", os.Args[0])
		fmt.Printf("For example: %s -dcinumber 12345 -cookie %s\n", os.Args[0], exampleCookie)
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
