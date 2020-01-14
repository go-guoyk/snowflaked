package main

type NextIDReq struct {
	Format string `query:"format"`
}

type NextIDRes struct {
	ID interface{} `json:"id"`
}

type NextIDsReq struct {
	Format string `query:"format"`
	Size   int    `query:"size"`
}

type NextIDsRes struct {
	IDs interface{} `json:"ids"`
}
