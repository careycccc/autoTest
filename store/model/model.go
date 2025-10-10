package model

// 基础结构体s
type BaseStruct struct {
	Random    any `json:"random"`
	Language  any `json:"language"`
	Signature any `json:"signature"`
	Timestamp any `json:"timestamp"`
}

// 请求头的基础结构体
type BaseHeaderStruct struct {
	Referer   any
	Origin    any
	Domainurl any
}

// 请求头的结构体 TenantId
type DeskHeaderTenantIdStruct struct {
	TenantId  any
	Referer   any
	Origin    any
	Domainurl any
}

// 请求头的结构体Authorization 和 TenantId
type DeskHeaderAuthorizationStruct struct {
	TenantId      any
	Referer       any
	Origin        any
	Domainurl     any
	Authorization any
}

// 后台请求结构体带Authorization
type AdminHeaderStruct struct {
	Referer       any
	Origin        any
	Domainurl     any
	Authorization any
}

// 查询相关的struct
type QueryPayloadStruct struct {
	PageNo    any `json:"pageNo"`
	PageSize  any `json:"pageSize"`
	OrderBy   any `json:"orderBy"`
	Random    any `json:"random"`
	Language  any `json:"language"`
	Signature any `json:"signature"`
	Timestamp any `json:"timestamp"`
}
