package api

type Attachment struct {
	BaseModel

	ID          string `json:"id"`
	ContentType string `json:"contentType"`
	Media       string `json:"media,omitempty"`
}

type ListAttachmentsResponse struct {
	Attachments []Attachment `json:"attachments"`
}
