package excel

type KVPair struct {
	Key   string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value int32  `protobuf:"varint,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

type KVPairs struct {
	Key    string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value1 int32  `protobuf:"varint,2,opt,name=Value1,proto3" json:"Value1,omitempty"`
	Value2 int32  `protobuf:"varint,3,opt,name=Value2,proto3" json:"Value2,omitempty"`
	Value3 int32  `protobuf:"varint,4,opt,name=Value3,proto3" json:"Value3,omitempty"`
	Value4 int32  `protobuf:"varint,5,opt,name=Value4,proto3" json:"Value4,omitempty"`
}

type KVPairInt struct {
	Key   int32 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value int32 `protobuf:"varint,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

type KVPairInt2 struct {
	Key    int32 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value1 int32 `protobuf:"varint,2,opt,name=Value1,proto3" json:"Value1,omitempty"`
	Value2 int32 `protobuf:"varint,3,opt,name=Value2,proto3" json:"Value2,omitempty"`
}

type TwoKVInt struct {
	Key1   int32 `protobuf:"varint,1,opt,name=Key1,proto3" json:"Key1,omitempty"`
	Key2   int32 `protobuf:"varint,2,opt,name=Key2,proto3" json:"Key2,omitempty"`
	Value1 int32 `protobuf:"varint,3,opt,name=Value1,proto3" json:"Value1,omitempty"`
	Value2 int32 `protobuf:"varint,4,opt,name=Value2,proto3" json:"Value2,omitempty"`
}

type SkillTriggerCondition struct {
	ParamGetterKey  string `protobuf:"bytes,1,opt,name=ParamGetterKey,proto3" json:"ParamGetterKey,omitempty"`
	ParamGetterArgs int32  `protobuf:"varint,2,opt,name=ParamGetterArgs,proto3" json:"ParamGetterArgs,omitempty"`
	CompareMode     string `protobuf:"bytes,3,opt,name=CompareMode,proto3" json:"CompareMode,omitempty"`
	CompareValue    int32  `protobuf:"varint,4,opt,name=CompareValue,proto3" json:"CompareValue,omitempty"`
}

type KVPairsInt struct {
	Key    int32 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value1 int32 `protobuf:"varint,2,opt,name=Value1,proto3" json:"Value1,omitempty"`
	Value2 int32 `protobuf:"varint,3,opt,name=Value2,proto3" json:"Value2,omitempty"`
	Value3 int32 `protobuf:"varint,4,opt,name=Value3,proto3" json:"Value3,omitempty"`
	Value4 int32 `protobuf:"varint,5,opt,name=Value4,proto3" json:"Value4,omitempty"`
	Value5 int32 `protobuf:"varint,6,opt,name=Value5,proto3" json:"Value5,omitempty"`
}

type MaxInt struct {
	Value int32 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

type KVPairString struct {
	Key   string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

type KVPairStrings struct {
	Key    string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value1 string `protobuf:"bytes,2,opt,name=Value1,proto3" json:"Value1,omitempty"`
	Value2 string `protobuf:"bytes,3,opt,name=Value2,proto3" json:"Value2,omitempty"`
	Value3 string `protobuf:"bytes,4,opt,name=Value3,proto3" json:"Value3,omitempty"`
	Value4 string `protobuf:"bytes,5,opt,name=Value4,proto3" json:"Value4,omitempty"`
	Value5 string `protobuf:"bytes,6,opt,name=Value5,proto3" json:"Value5,omitempty"`
}
