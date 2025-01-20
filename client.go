package blastenginego

type Client struct {
	apiKey string
	userId string
}

func initializeClient(apiKey string, userId string) Client {
	// Initialize the client
	c := Client{apiKey: apiKey, userId: userId}
	return c
}
