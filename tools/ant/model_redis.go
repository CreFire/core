// 模型来自pb
// //特别注意，lua只至此double，int64的数据如果进行cmsgpack打包解包可能出现精度问题导致bug
package antnet

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/vmihailenco/msgpack"
	"reflect"
)

type RedisModel struct{}

func (r *RedisModel) DBData(v proto.Message) []byte {
	return DBData(v)
}

func (r *RedisModel) DBStr(v proto.Message) string {
	return DBStr(v)
}

func (r *RedisModel) PbData(v proto.Message) []byte {
	return PbData(v)
}

func (r *RedisModel) PbStr(v proto.Message) string {
	return PbStr(v)
}

func (r *RedisModel) ParseDBData(data []byte, v proto.Message) bool {
	return ParseDBData(data, v)
}

func (r *RedisModel) ParseDBStr(str string, v proto.Message) bool {
	return ParseDBStr(str, v)
}

func (r *RedisModel) ParsePbData(data []byte, v proto.Message) bool {
	return ParsePbData(data, v)
}

func (r *RedisModel) ParsePbStr(str string, v proto.Message) bool {
	return ParsePbStr(str, v)
}

func DBDataAutoCompress(v proto.Message) []byte {
	if data, err := msgpack.Marshal(v); err == nil {
		dLen := len(data)
		if dLen > 1500 {
			compressedData := GZipCompress(data)
			text := make([]byte, len(compressedData)+1)
			text[0] = 1
			copy(text[1:], compressedData)
			return text
		}
		return data
	} else {
		panic(Sprintf("DBDataAutoCompress error msg: %v, err: %v", v, err.Error()))
	}
}

func DBStrAutoCompress(v proto.Message) string {
	return string(DBDataAutoCompress(v))
}

func ParseDBDataAutoUnCompress(data []byte, v proto.Message) bool {
	if data[0] == 1 {
		text, err := GZipUnCompress(data[1:])
		if err != nil {
			LogError("ParseDBDataAutoUnCompress GZipUnCompress error: %v, %v", v, err)
			return false
		}
		err = msgpack.Unmarshal(text, v)
		if err != nil {
			LogError("ParseDBDataAutoUnCompress Unmarshal error: %v, %v", v, err)
			return false
		}
	} else {
		err := msgpack.Unmarshal(data[1:], v)
		if err != nil {
			LogError("ParseDBDataAutoUnCompress Unmarshal error: %v, %v", v, err)
			return false
		}
	}
	return true
}

func ParseDBStrAutoUnCompress(str string, v proto.Message) bool {
	return ParseDBDataAutoUnCompress([]byte(str), v)
}

func DBData(v proto.Message) []byte {
	if data, err := msgpack.Marshal(v); err == nil {
		return data
	} else {
		panic(Sprintf("DBData error msg: %v, err: %v", v, err.Error()))
	}
}

func DBStr(v interface{}) string {
	if data, err := msgpack.Marshal(v); err == nil {
		return string(data)
	} else {
		panic(Sprintf("DBStr error msg: %v, err: %v", v, err.Error()))
	}
}

func PbData(v proto.Message) []byte {
	if data, err := proto.Marshal(v); err == nil {
		return data
	} else {
		panic(Sprintf("PbData error msg: %v, type: %v, err: %v", v, reflect.TypeOf(v), err.Error()))
	}
}

func PbStr(v proto.Message) string {
	if data, err := proto.Marshal(v); err == nil {
		return string(data)
	} else {
		panic(Sprintf("PbStr error msg: %v, err: %v", v, err.Error()))
	}
}

func ParseDBData(data []byte, v interface{}) bool {
	err := msgpack.Unmarshal(data, v)
	if err != nil {
		LogError("msgpack:%s struct:%v parse dbdata failed:%s", data, v, err.Error())
	}
	return err == nil
}

func ParseDBStr(str string, v interface{}) bool {
	err := msgpack.Unmarshal([]byte(str), v)
	if err != nil {
		LogError("msgpack:%s struct:%s parse dbstr failed:%s", str, v, err.Error())
	}
	return err == nil
}

func ParsePbData(data []byte, v proto.Message) bool {
	err := proto.Unmarshal(data, v)
	if err != nil {
		LogError("ParsePbData:%s struct:%v parse dbdata failed:%s", data, v, err.Error())
	}
	return err == nil
}

func ParsePbStr(str string, v proto.Message) bool {
	err := proto.Unmarshal([]byte(str), v)
	if err != nil {
		LogError("ParsePbStr:%s struct:%s parse pbstr failed:%s", str, v, err.Error())
	}
	return err == nil
}

func JsonStr(v interface{}) (string, error) {
	if data, err := json.Marshal(v); err == nil {
		return string(data), nil
	} else {
		return "", err
	}
}

func JsonData(v interface{}) ([]byte, error) {
	if data, err := json.Marshal(v); err == nil {
		return data, nil
	} else {
		return nil, err
	}
}

func ParseJsonData(data []byte, v interface{}) bool {
	err := json.Unmarshal(data, v)
	if err != nil {
		LogError("json:%s struct:%v parse jsondata failed:%s", data, v, err.Error())
	}
	return err == nil
}

func ParseJsonStr(str string, v interface{}) error {
	return json.Unmarshal([]byte(str), v)
}
