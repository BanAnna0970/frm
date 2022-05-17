# fm

```go
func getReactInfo(authorization, channel, proxy string) (string, string) {
	fmt.Println(authorization, channel)
	headers := make(map[string]string)
	headers["Host"] = "discord.com"
	headers["x-discord-locale"] = "en-GB"
	headers["x-debug-options"] = "bugReporterEnabled"
	headers["accept-language"] = "en-US,en-RU;q=0.9,ru-RU;q=0.8"
	headers["authorization"] = authorization
	headers["user-agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) discord/0.0.266 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36"
	headers["accept"] = "*/*"
	headers["sec-fetch-site"] = "same-origin"
	headers["sec-fetch-mode"] = "cors"
	headers["sec-fetch-dest"] = "empty"

	data := req.FastData{
		Headers: headers,
		URL:     fmt.Sprintf("https://discord.com/api/v9/channels/%v/messages", channel),
		Method:  "GET",
	}

	fh := data.Build()

	fh.Client.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
	err := fh.DoRequest()
	if err != nil {
		logger.Logger.Error().Err(err)
	}

	body := fh.Request.Body()
	authorID, _, _, err := jsonparser.Get(body, "[0]", "author", "id")
	if err != nil {
		logger.Logger.Error().Err(err)
	}
	customID, _, _, err := jsonparser.Get(body, "[0]", "components", "[0]", "components", "[0]", "custom_id")
	if err != nil {
		logger.Logger.Error().Err(err)
	}

	return string(authorID), string(customID)
}
