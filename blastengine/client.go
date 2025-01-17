package blastengine

type Client struct {
    UserID string
    APIKey string
}

func Initialize(userID, apiKey string) Client {
    return Client{
        UserID: userID,
        APIKey: apiKey,
    }
}
