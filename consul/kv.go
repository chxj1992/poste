package consul

import (
	"github.com/hashicorp/consul/api"
	"poste/util"
)

func KVSet(key string, value string) bool {
	p := &api.KVPair{Key: key, Value: []byte(value)}
	kv := GetKV()
	_, err := kv.Put(p, nil)
	if err != nil {
		util.LogError("consul kv set value failed. error : %s", err)
		return false
	}
	return true
}

func KVGet(key string) string {
	kv := GetKV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		util.LogError("consul kv get value failed. error : %s", err)
		return ""
	}
	if pair == nil {
		return ""
	}
	return string(pair.Value)
}

func KVDelete(key string) bool {
	_, err := GetKV().Delete(key, nil)
	if err != nil {
		util.LogError("consul kv delete value failed. error : %s", err)
		return false
	}
	return true
}

func KVClear() bool {
	_, err := GetKV().DeleteTree("", nil)
	if err != nil {
		util.LogError("consul kv clear failed. error : %s", err)
		return false
	}
	return true
}

func KVValues(prefix string) []string {
	pairs, _, err := GetKV().List(prefix, nil)

	values := []string{}
	if err != nil {
		util.LogError("consul kv get value failed. error : %s", err)
		return values
	}
	if pairs == nil {
		return values
	}
	for _, pair := range pairs {
		values = append(values, string(pair.Value))
	}
	return values
}

func KVWatch(prefix string, callback func(values []*api.KVPair)) {
	q := &api.QueryOptions{}
	for {
		pairs, meta, err := GetKV().List(prefix, q)
		q.WaitIndex = meta.LastIndex

		values := []*api.KVPair{}
		if err != nil {
			util.LogError("consul kv get value failed. error : %s", err)
		}
		if pairs != nil {
			for _, pair := range pairs {
				values = append(values, pair)
			}
		}
		callback(values)
	}
}
