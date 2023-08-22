package antnet

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"sort"
	"strings"
)

//该方法会绕过压缩和加密
//func Broadcast(msg *Message, fun func(msgque IMsgQue) bool) {
//	if msg == nil {
//		return
//	}
//	c := make(chan struct{})
//	gmsgMapSync.Lock()
//	gmsg := gmsgArray[gmsgId]
//	gmsgArray[gmsgId+1] = &gMsg{c: c}
//	gmsgId++
//	gmsgMapSync.Unlock()
//	gmsg.msg = msg
//	gmsg.fun = fun
//	close(gmsg.c)
//}

func HttpGetWithBasicAuth(url, name, passwd string) (string, error, *http.Response) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", ErrHttpRequest, nil
	}
	req.SetBasicAuth(name, passwd)
	resp, err := client.Do(req)
	if err != nil {
		return "", ErrHttpRequest, nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", ErrHttpRequest, nil
	}
	resp.Body.Close()
	return string(body), nil, resp
}

func HttpGet(url string) (string, error, *http.Response) {
	resp, err := http.Get(url)
	if err != nil {
		return "", ErrHttpRequest, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", ErrHttpRequest, resp
	}
	resp.Body.Close()
	return string(body), nil, resp
}

func HttpPostPostForm(url string, form url.Values) (string, error, *http.Response) {
	return HttpPost(url, form.Encode())
}

func HttpPost(url, form string) (string, error, *http.Response) {
	client := &http.Client{}
	resp, err := client.Post(url, "application/x-www-form-urlencoded", strings.NewReader(form))
	if err != nil {
		return "", err, nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err, resp
	}
	return string(body), nil, resp

	//以下存在并发问题 http2: client connection force closed via ClientConn.Close
	//resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(form))
	//if err != nil {
	//	return "", err, nil
	//}
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return "", err, resp
	//}
	//resp.Body.Close()
	//return string(body), nil, resp
	//http.DefaultClient.CloseIdleConnections()
}

func PostFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		LogError("error writing to buffer")
		return err
	}
	//打开文件句柄操作
	//fmt.Println("filename", filename)
	fh, err := os.Open(filename)
	if err != nil {
		LogError("error opening file")
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//resp_body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return err
	//}
	//fmt.Println(resp.Status)
	//fmt.Println(string(resp_body))
	//fmt.Println("filename", filename)

	return nil
}

func HttpUpload(url, field, file string) (*http.Response, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile(field, file)
	if err != nil {
		LogError("create form file failed:%s\n", err)
		return nil, err
	}

	srcFile, err := os.Open(file)
	if err != nil {
		LogError("%open source file failed:%s\n", err)
		return nil, err
	}
	defer srcFile.Close()
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		LogError("write to form file falied:%s\n", err)
		return nil, err
	}

	contentType := writer.FormDataContentType()
	writer.Close()
	resp, err := http.Post(url, contentType, buf)
	if err != nil {
		LogError("post failed:%s\n", err)
	}

	return resp, err
}

func SendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

var allIp []net.IP

var excludeNetName = []string{"aioCloud"}

func inExcludeNetName(name string) bool {
	for _, v := range excludeNetName {
		if v == name {
			return true
		}
	}
	return false
}

func GetSelfIp(ifnames ...string) []net.IP {
	if allIp != nil {
		return allIp
	}
	inters, _ := net.Interfaces()
	if len(ifnames) == 0 {
		ifnames = []string{"eth", "eno", "lo", "ens33", "无线网络连接", "本地连接", "以太网"}
	}
	filterFunc := func(name string) bool {
		for _, v := range ifnames {
			if strings.Index(name, v) != -1 {
				return true
			}
		}
		return false
	}

	for _, inter := range inters {
		if inExcludeNetName(inter.Name) {
			continue
		}
		if !filterFunc(inter.Name) {
			continue
		}
		LogInfo("Network Adapter:%s", inter.Name)
		addrs, _ := inter.Addrs()
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() {
					allIp = append(allIp, ipnet.IP)
					LogInfo("IP:%s", ipnet.IP.To4().String())
				}
			}
		}
	}
	if len(allIp) == 0 {
		panic(basal.NewError("找不到网卡地址"))
	}
	return allIp
}

func GetSelfIntraIp(ifnames ...string) (ips []string) {
	all := GetSelfIp(ifnames...)
	for _, v := range all {
		if v.IsLoopback() || v.IsLinkLocalMulticast() || v.IsLinkLocalUnicast() {
			continue
		}
		if !IsExtraIP(v) {
			ips = append(ips, v.To4().String())
		}
	}
	sort.Slice(ips, func(i, j int) bool {
		vi := Atoi32(strings.Split(ips[i], ".")[0])
		vj := Atoi32(strings.Split(ips[j], ".")[0])
		return vi > vj
	})

	return
}

func GetSelfExtraIp(ifnames ...string) (ips []string) {
	all := GetSelfIp(ifnames...)
	for _, v := range all {
		if IsExtraIP(v) {
			ips = append(ips, v.To4().String())
		}
	}
	return
}

func IsExtraIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

func IPCanUse(ip string) bool {
	var err error
	for port := 1024; port < 65535; port++ {
		addr := Sprintf("%v:%v", ip, port)
		listen, err := net.Listen("tcp", addr)
		if err == nil {
			listen.Close()
			break
		} else if StrFind(err.Error(), "address is not valid") != -1 {
			return false
		}
	}
	return err == nil
}
