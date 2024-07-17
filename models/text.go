package models

type Text []byte

func (text *Text) String() string {
	return string(*text)
}
