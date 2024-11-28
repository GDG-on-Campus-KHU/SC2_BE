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

// 재난 응답 구조체
type DisasterGuideResponse struct {
	PushAlarming   string         `json:"push_alarming"`
	Congestion     CongestionInfo `json:"congestion"`
	DisasterRadius float64        `json:"disaster_radius"`
	ActionPlan     []ActionPlan   `json:"action_plan"`
}

type CongestionInfo struct {
	AreaName        string `json:"area_nm"`
	LivePopulation  *int   `json:"live_ppltn_stts"` // Null 값을 허용하기 위해 포인터 사용
	CongestionLevel string `json:"area_congest_lvl"`
	CongestionMsg   string `json:"area_congest_msg"`
	PopulationMin   string `json:"area_ppltn_min"`
	PopulationMax   string `json:"area_ppltn_max"`
	WarnValue       string `json:"warn_val"`
	WarnStress      string `json:"warn_stress"`
	AnnounceTime    string `json:"announce_time"`
	Command         string `json:"command"`
	CancelYN        string `json:"cancel_yn"`
	WarnMsg         string `json:"warn_msg"`
}

// 행동 요령 구조체 정의
type ActionPlan struct {
	SafetyCategoryName1 string  `json:"safety_cate_nm1"`
	SafetyCategoryName2 string  `json:"safety_cate_nm2"`
	SafetyCategoryName3 string  `json:"safety_cate_nm3"`
	SafetyCategoryName4 *string `json:"safety_cate_nm4"`
	ActionRemarks       string  `json:"actRmks"`
	Category1           string  `json:"safety_cate1"`
	Category2           string  `json:"safety_cate2"`
	Category3           string  `json:"safety_cate3"`
	ContentsURL         *string `json:"contentsUrl"`
}
