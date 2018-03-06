package main

type TemplateData struct {
	First    DataItem `json:"first"`
	Keyword1 DataItem `json:"keyword1"`
	Keyword2 DataItem `json:"keyword2"`
	Remark   DataItem `json:"remark"`
}

type DataItem struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}
