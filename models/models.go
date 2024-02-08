package models

type Expression struct {
	Expression   string `json:"expression"`
	AddTime      int64  `json:"addTime"`
	SubTime      int64  `json:"subTime"`
	MultiplyTime int64  `json:"multiplyTime"`
	DivideTime   int64  `json:"divideTime"`
}
