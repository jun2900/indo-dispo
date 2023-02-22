package entity

type Attachment struct {
	Name string `json:"attachment_name"`
	File []byte `json:"attachment_file"`
}
