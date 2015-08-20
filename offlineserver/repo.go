package main

import (
	"log"
)

func RepoCreateATSRsp(req *ATSReq) []ATSRsp {
	var rest []ATSRsp

	for _, atsR := range req.SkuIds {
		rsp := ATSRsp{
			SkuId:          atsR,
			Ats:            10,
			AllowBackOrder: true,
		}
		rest = append(rest, rsp)
	}

	log.Printf("ATS Rsp %+v\n", rest)

	return rest
}
