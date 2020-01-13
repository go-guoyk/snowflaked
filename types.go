package main

type CreateResp struct {
	ID int64 `json:"id"`
}

type CreateSResp struct {
	ID string `json:"id"`
}

type BatchReq struct {
	Size int `json:"size"`
}

type BatchResp struct {
	IDs []int64 `json:"ids"`
}

type BatchSResp struct {
	IDs []string `json:"ids"`
}
