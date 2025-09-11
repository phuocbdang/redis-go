package core

import (
	"errors"
	"redisgo/internal/data_structure"
)

func cmdSADD(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SADD' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		set = data_structure.NewSimpleSet(key)
		setStore[key] = set
	}
	count := set.Add(args[1:]...)
	return Encode(count, false)
}

func cmdSREM(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SREM' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		return Encode(0, false)
	}
	count := set.Rem(args[1:]...)
	return Encode(count, false)
}

func cmdSMEMBERS(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SMEMBERS' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		return Encode(make([]string, 0), false)
	}
	members := set.Members()
	return Encode(members, false)
}

func cmdSISMEMBER(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SISMEMBER' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		return Encode(0, false)
	}
	isMember := set.IsMember(args[1])
	return Encode(isMember, false)
}
