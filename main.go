package main

import "fmt"

type IdentifyByte struct {
	Offset    int    `json:"offset"`
	TypeBytes string `json:"typeBytes"`
}

type FileTypeAllow struct {
	IdentifyBytes []IdentifyByte `json:"IdentifyBytes"`
	TypeName      string         `json:"typeName"`
}

type FilterPolicy struct {
	ContentBlackList []string `json:"contentBlackList"`
	//ContentWhiteList []string        `json:"contentWhiteList"`
	FileTypeAllows []FileTypeAllow `json:"fileTypeAllow"`
	SizeMaxKb      int             `json:"sizeMaxKB"`
	SizeMinKb      int             `json:"sizeMinKB"`
	SubFixAllow    string          `json:"subFixAllow"`
}

type Policy struct {
	FtpCmd         string         `json:"ftpCmd"`
	FtpUsers       []string       `json:"ftpUsers"`
	Password       string         `json:"password"`
	Protocol       string         `json:"protocol"`
	RefreshSeconds int            `json:"refreshSeconds"`
	RemotePath     string         `json:"remotePath"`
	ServerIp       string         `json:"serverIp"`
	ServerPort     string         `json:"serverPort"`
	Username       string         `json:"username"`
	CacheDay       int            `json:"cacheDay"`
	FilterPolicys  []FilterPolicy `json:"filterPolicy"`
}

func main() {
	//var policys []Policy
	//s := "[{\"cacheDays\":1,\"filterPolicy\":[{\"contentBlackList\":[],\"contentWhiteList\":[{\"regx\":\"[0-9]\"}],\"fileTypeAllow\":[{\"IdentifyBytes\":[{\"offset\":0,\"typeBytes\":\"D0CF11E0\"}],\"typeName\":\"doc\"},{\"IdentifyBytes\":[{\"offset\":0,\"typeBytes\":\"UTF-8,Unicode,ANSI\"},{\"offset\":0,\"typeBytes\":\"UTF-8,Unicode,ANSI\"},{\"offset\":0,\"typeBytes\":\"UTF-8,Unicode,ANSI\"}],\"typeName\":\"txt\"}],\"sizeMaxKB\":9999999,\"sizeMinKB\":0,\"subFixAllow\":\".txt\"},{\"contentBlackList\":[],\"contentWhiteList\":[{\"$ref\":\"$.data.netFilePolicy.Policies[0].filterPolicy[0].contentWhiteList[0]\"}],\"fileTypeAllow\":[{\"IdentifyBytes\":[{\"offset\":0,\"typeBytes\":\"504B0304\"}],\"typeName\":\"zip\"}],\"sizeMaxKB\":10000,\"sizeMinKB\":0,\"subFixAllow\":\"zip;\"}],\"ftpCmd\":\"d\",\"ftpUsers\":[],\"password\":\"123\",\"protocol\":\"FTP\",\"refreshSeconds\":12,\"remote_path\":\"/\",\"serverIp\":\"192.2.2.1\",\"serverPort\":\"90\",\"username\":\"123\"}]"
	//err := json.Unmarshal([]byte(s), &policys)
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Printf("%+v\n", policys)
	//}
	min := 5
	max := 5
	if min == 0 && max == 0 || max < min {
		// 非法参数直接通过
		fmt.Println("===")
	}

}
