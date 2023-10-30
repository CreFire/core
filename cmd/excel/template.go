package excel

import "fmt"

var template = `
fname = "%s"
filename = filepath.Join(confDir, basal.Sprintf("%%v.bytes", fname))
if data, err := os.ReadFile(filename); err != nil {
antnet.LogError("docpb file load failed: %%v, %%v, %%v", confDir, err, filename)
return err
} else {
pbcfg := &%sCfg{}
err = antnet.PBUnPack(data, pbcfg)
if err != nil {
antnet.LogError("docpb file unpack failed: %%v, %%v, %%v", confDir, err, filename)
return err
}
m.%s = make(map[int32]*%sData, len(pbcfg.Configs))
jsOld := make([]*pb.DocConfigDataID, 0, len(pbcfg.Configs))
jsFileName := filepath.Join(toJsDir, basal.Sprintf("%%v.json", fname))
err = jsonex.LoadFileTo(jsFileName, &jsOld)
if err != nil && !os.IsNotExist(err) {
antnet.LogError("jsonex.LoadFileTo err: %%v, %%v", jsFileName, err)
return err
}
jsNew := make([]*%sData, 0, len(pbcfg.Configs))
for _, conf := range pbcfg.Configs {
if conf.GetID() == docPbInt32Null {
continue
}
if _, ok := m.%s[conf.GetID()]; ok{
err = basal.NewError("配置档ID重复: %%v, %%v", filename, conf.GetID())
antnet.LogError(err)
return err
}
m.%s[conf.GetID()] = conf
jsNew = append(jsNew, conf)
}

for _, v := range jsOld {
err = checkOldAndNewConfigData(v, m.%s[v.GetID()])
if err != nil {
antnet.LogError("checkOldAndNewConfigData err: %%v, %%v", filename, err)
return err
}
}

js, err := jsonex.TryDump(jsNew, true)
if err != nil {
antnet.LogError("ToJsonString err: %%v, %%v", filename, err)
return err
}

jsFile, err := basal.OpenFileB(jsFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
if err != nil {
antnet.LogError("docpb file convert to json failed: %%v, %%v, %%v", confDir, err, jsFileName)
return err
}

_, err = jsFile.WriteString(js)
if err != nil {
antnet.LogError("docpb WriteString failed: %%v, %%v, %%v", confDir, err, jsFileName)
return err
}
//err = jsFile.Sync()
//if err != nil {
//    antnet.LogError("docpb Sync failed: %%v, %%v, %%v", confDir, err, jsFileName)
//    return err
//}
err = jsFile.Close()
if err != nil {
antnet.LogError("docpb Close failed: %%v, %%v, %%v", confDir, err, jsFileName)
return err
}
}`

func Generate(fname string, modelnames ...string) {

	_ = fmt.Sprintf(template, fname, modelnames)
}
