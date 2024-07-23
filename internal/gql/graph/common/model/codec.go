package model

import "encoding/json"

func (a ShellConfiguration) MarshalBinary() (data []byte, err error) {
	return json.Marshal(a)
}

func (a *ShellConfiguration) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &a)
}
