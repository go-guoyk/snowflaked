package main

type CreateReq struct {
	Size int `json:"size"`
}

type CreateResp struct {
	IDs []int64 `json:"ids"`
}
