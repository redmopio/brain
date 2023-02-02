package channels

type TelegramConnector struct {
	apiKey string
}

func NewTelegramConnector(apiKey string) *TelegramConnector {
	return &TelegramConnector{
		apiKey: apiKey,
	}
}
