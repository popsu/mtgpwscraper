# Magic the Gathering Planeswalker Points history scraper

## Requirements:
- Your DCI Number

## Additional requirements for detailed event history:
- Login access to your Wizards account that is linked to your DCI number

## Usage

1. Download the binary from releases.
2. Login to your Wizards account at https://www.wizards.com/Magic/PlaneswalkerPoints/
3. While logged in, click CTRL + Shift + J in your browser to open developer console. Go to Application tab -> Storage -> Cookies -> `https://www.wizards.com/` and from the list on the right find the cookie Named PWP.ASPXAUTH and copy paste the Value from next to it. Pass this value to this program as the -cookie argument.
4. Pass your dcinumber and the cookie you just got to this program in following format:
` pwpscraper.exe -dcinumber $DCINUMBER -cookie $COOKIE for example:
```
pwpscraper.exe -dcinumber 12345 -cookie 9A3D5E64320129F4764627A2E03F6025B4F838A2D40C5FE255809E2CC6F18EBCFAF0CA38A986460B7075183ACBF2385468C11F943BA4E4E2648B640F45C9B4E33643C150E664D52897F6DC35FC6DEBA0B3979E67497B3029FC48058388590A2C986FEC3F0894ACA1C90190064879ECEB8E836295803B4E25EB3C8A0EF0091090AAEF60A98463B8C68C0E7BBFE27958F9F2FE424474D752F1A413CCC3ABD8558CD3449EF92446C7733B23C1E9ABBB5FF240261F9D1579A910E074AF494A51C82D79399E158533F2697A2DC7286FF56B2636D6FAF67EBEFD983924DB1646487E2189F683DF2361A94B80EAE6A2B6C03A48FEC9ECCCD5DF200C7D2F08A606A3297B33B514942910D55D0D638D7C690FE18883331EBDC00F6F8CEF733720EFA46E5D5D1335205211EDDCF7952EDECC52183A5A66D3B7381149ABB32F467B985A547EB0E3A5DA3AFC26B0E0BDA548763C082FCA5523AC92DBCAF7A72D234E755B5F36989C714A87F9F4FAEBC0030384150780BC6758AB697637D4E7B318E6029F6FA183CCF8D62DD9B740F16248935C513484CBA91323EFCF61F7AE9E884B5FBA96A61FD87366A31C09473181F74637D06385B31BC390
```

- The program will go and download your event history and save the detailed html's in`eventdata` folder.
- The program will parse the downloaded eventdata and combine them into `planeswalker.json` file that contains all the information. You can delete the `eventdata` folder afterwards.
- **If you passed invalid cookie value (or if you don't have Wizards account at all you can pass fake cookie to get your history without details) the final `planeswalker.json` contains empty EventDetails objects.**
