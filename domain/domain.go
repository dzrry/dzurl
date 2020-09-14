package domain

type Redirect struct {
	Key       string `json:"key" msgpack:"key" valid:"-"`
	URL       string `json:"url" msgpack:"url" valid:"requrl"`
	CreatedAt int64  `json:"created_at" msgpack:"created_at" valid:"-"`
}
