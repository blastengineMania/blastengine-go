package blastengine

func initialize(apiKey string, userId string) Client {
	// Initialize the client
	c := Client{apiKey: apiKey, userId: userId}
	return c
}
