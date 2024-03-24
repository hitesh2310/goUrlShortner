package models

type Request struct {
	Url string `json:"url"`
}

type RedisEntry struct {
	LongUrl   string `json:"longUrl"`
	EpochTime int    `json:"epochTime"`
}
