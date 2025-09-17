package core

import (
	"errors"
	"fmt"
	"math"
	"redisgo/internal/constant"
	"redisgo/internal/data_structure"
	"strconv"
)

func cmdCMSINITBYDIM(args []string) []byte {
	if len(args) != 3 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'CMS.INITBYDIM' command"), false)
	}
	key := args[0]
	width, err := strconv.ParseUint(args[1], 10, 32)
	if err != nil {
		return Encode(fmt.Errorf("width must be a integer number %s", args[1]), false)
	}
	depth, err := strconv.ParseUint(args[2], 10, 32)
	if err != nil {
		return Encode(fmt.Errorf("depth must be a integer number %s", args[2]), false)
	}
	_, exist := cmsStore[key]
	if exist {
		return Encode(errors.New("CMS: key already exist"), false)
	}
	cmsStore[key] = data_structure.CreateCMS(uint32(width), uint32(depth))
	return constant.RespOk
}

func cmdCMSINITBYPROB(args []string) []byte {
	if len(args) != 3 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'CMS.INITBYPROB' command"), false)
	}
	key := args[0]
	errRate, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return Encode(fmt.Errorf("errRate must be a floating point number %s", args[1]), false)
	}
	if errRate >= 1 || errRate <= 0 {
		return Encode(errors.New("CMS: invalid overestimation value"), false)
	}
	probability, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return Encode(fmt.Errorf("probability must be a floating point number %s", args[2]), false)
	}
	if probability >= 1 || probability <= 0 {
		return Encode(errors.New("CMS: invalid prob value"), false)
	}
	_, exist := cmsStore[key]
	if exist {
		return Encode(errors.New("CMS: key already exist"), false)
	}
	w, d := data_structure.CalcCMSDim(errRate, probability)
	cmsStore[key] = data_structure.CreateCMS(w, d)
	return constant.RespOk
}

func cmdCMSINCRBY(args []string) []byte {
	if len(args) < 3 || len(args)%2 == 0 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'CMS.INCRBY' command"), false)
	}
	key := args[0]
	cms, exist := cmsStore[key]
	if !exist {
		return Encode(errors.New("CMS: key does not exist"), false)
	}
	var res []interface{}
	for i := 1; i < len(args); i += 2 {
		item := args[i]
		value, err := strconv.ParseUint(args[i+1], 10, 32)
		if err != nil {
			return Encode(fmt.Errorf("increase must be a non negative integer number %s", args[i+1]), false)
		}
		count := cms.IncrBy(item, uint32(value))
		if count == math.MaxUint32 {
			res = append(res, "CMS: INCRBY overflow")
		} else {
			res = append(res, int64(count))
		}
	}
	return Encode(res, false)
}

func cmdCMSQUERY(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'CMS.QUERY' command"), false)
	}
	key := args[0]
	cms, exist := cmsStore[key]
	if !exist {
		return Encode(errors.New("CMS: key does not exist"), false)
	}
	var res []int64
	for i := 1; i < len(args); i++ {
		item := args[i]
		res = append(res, int64(cms.Count(item)))
	}
	return Encode(res, false)
}
