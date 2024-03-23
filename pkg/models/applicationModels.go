package models

type Request struct {
	Url         string `json:"url"`
	ShortString string `json:"shortString"`
}

type RedisEntry struct {
	LongUrl   string `json:"longUrl"`
	EpochTime int    `json:"epochTime"`
}
