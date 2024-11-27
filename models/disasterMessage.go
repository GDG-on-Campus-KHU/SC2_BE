package models

// DisasterMessage 구조체 정의
type DisasterMessage struct {
	SN           string `json:"SN"`
	CRT_DT       string `json:"CRT_DT"`
	MSG_CN       string `json:"MSG_CN"`
	RCPTN_RGN_NM string `json:"RCPTN_RGN_NM"`
	EMRG_STEP_NM string `json:"EMRG_STEP_NM"`
	DST_SE_NM    string `json:"DST_SE_NM"`
	REG_YMD      string `json:"REG_YMD"`
	MDFCN_YMD    string `json:"MDFCN_YMD"`
}

// DisasterResponse 구조체 정의
type DisasterResponse struct {
	ResponseCode string            `json:"responseCode"`
	ResponseMsg  string            `json:"responseMsg"`
	Items        []DisasterMessage `json:"items"`
}
