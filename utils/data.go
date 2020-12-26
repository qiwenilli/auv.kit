package utils

import (
	"errors"
	"reflect"

	"github.com/gogf/gf/util/gconv"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/twitchtv/twirp"
)

func ToProtobufStruct(data interface{}, time2str bool) (*structpb.Struct, error) {
	//
	dataMap := make(map[string]interface{})
	//
	typof := reflect.TypeOf(data)
	valof := reflect.ValueOf(data)
	for i := 0; i < valof.NumField(); i++ {
		typ := typof.Field(i)
		val := valof.Field(i)

		if typ.Type.String() == "time.Time" {
			// 转成时间戳
			timeInf := gconv.Convert(val.Interface(), typ.Type.String())
			// timeVal := gconv.Time(timeInf, "2006-01-02 15:04:05")
			timeVal := gconv.Time(timeInf)
			if time2str {
				dataMap[typ.Name] = timeVal.Format("2006-01-02 15:04:05")

			} else {
				dataMap[typ.Name] = timeVal.UnixNano()

			}

		} else if val.Kind().String() == "struct" {
			return nil, errors.New("element is struct")

		} else {
			dataMap[typ.Name] = val.Interface()

		}

	}
	return structpb.NewStruct(dataMap)

}

func ToTwirpError(err error) error {
	if err == nil {
		return nil
	}
	return twirp.NewError(twirp.NoError, err.Error())
}
