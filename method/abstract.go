package method

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"webdav-aliyundriver/config"
	"webdav-aliyundriver/model"
)

const (
	Iso8601        = "yyyy-MM-dd'T'HH:mm:ss'Z'"
	Rfc1123        = "EEE, dd MMM yyyy HH:mm:ss z"
	DDMMYYHHMMSS   = "dd/MM/yy' 'HH:mm:ss"
	DefaultTimeout = 3600
	Infinity       = 3
)
const (
	DestinationKey = "Destination"
)

func ParseDestinationHeader(response http.ResponseWriter, request *http.Request) (string, error) {
	// 获取文件路径
	path := request.Header.Get(DestinationKey)
	if len(path) <= 0 {
		response.WriteHeader(http.StatusBadRequest)
		return "", nil
	}
	// Remove url encoding from destination
	destinationPath, err := url.PathUnescape(path)
	if err != nil {
		destinationPath = path
	}
	protocolIndex := strings.Index(destinationPath, "://")
	if protocolIndex >= 0 {
		//如果目标URL包含协议，我们可以安全地 将“://”后的第一个“/”字符进行修剪 ://xxx -> / ://xxx/bb -> /bb
		firstSeparator := strings.Index(destinationPath[protocolIndex+4:], "/")
		if firstSeparator < 0 {
			destinationPath = "/"
		} else {
			destinationPath = destinationPath[firstSeparator:]
		}
	} else {
		hostName := request.Host
		if len(hostName) > 0 && strings.HasPrefix(destinationPath, hostName) {
			//直接获取路径 需要通过host对比 host:ggg  destinationPath:ggg/sss -> /sss
			destinationPath = destinationPath[len(hostName):]
		}
		portIndex := strings.Index(destinationPath, ":")
		if portIndex >= 0 {
			destinationPath = destinationPath[portIndex:]
		}
		if strings.HasPrefix(destinationPath, ":") {
			firstSeparator := strings.Index(destinationPath, "/")
			if firstSeparator < 0 {
				destinationPath = "/"
			} else {
				destinationPath = destinationPath[firstSeparator:]
			}
		}
	}
	destinationPath = Normalize(destinationPath)
	if len(config.WebConf.ContextPath) > 0 && strings.HasPrefix(destinationPath, config.WebConf.ContextPath) {
		destinationPath = destinationPath[len(config.WebConf.ContextPath):]
	}
	return destinationPath, nil
}

//Normalize 返回一个上下文相关的路径，以"/"开头，表示解析出".."和"."元素后指定路径的规范版本。
//如果指定的路径试图超出当前上下文的边界(即存在太多的".."路径元素)，则返回<code>null</code>。
func Normalize(path string) string {
	if len(path) <= 0 {
		return ""
	}
	normalized := path
	if normalized == "/." {
		return "/"
	}
	if strings.Index(normalized, "\\") >= 0 {
		normalized = strings.Replace(normalized, "\\", "/", -1)
		if strings.HasPrefix(normalized, "/") {
			normalized = "/" + normalized
		}
	}
	// 解析在规范化路径中出现的"//"
	for {
		index := strings.Index(normalized, "//")
		if index < 0 {
			break
		}
		normalized = normalized[0:index] + normalized[index+1:]
	}
	// 解决“/”的出现。 在规范化路径中
	for {
		index := strings.Index(normalized, "/./")
		if index < 0 {
			break
		}
		normalized = normalized[0:index] + normalized[index+2:]
	}
	// 解决“/.. 在规范化路径中
	for {
		index := strings.Index(normalized, "/../")
		if index < 0 {
			break
		}
		if index == 0 {
			return ""
		}
		index2 := strings.LastIndex(normalized[index-1:], "/")
		normalized = normalized[0:index2] + normalized[index+3:]
	}
	// 返回我们已经完成的规范化路径
	return normalized
}

func RelativePath(r *http.Request) string {
	destinationPath := r.RequestURI
	if len(config.WebConf.ContextPath) > 0 && strings.HasPrefix(destinationPath, config.WebConf.ContextPath) {
		destinationPath = destinationPath[len(config.WebConf.ContextPath):]
	}
	return destinationPath
}

//ParentPath 通过删除最后一个'/'及其之后的所有内容，从给定路径创建父路径
func ParentPath(path string) string {
	index := strings.LastIndex(path, "/")
	if index != -1 {
		return path[:index]
	}
	return ""
}

//CleanPath 移除路径字符串末尾的/(如果存在的话)
func CleanPath(path string) string {
	if strings.HasSuffix(path, "/") && len(path) > 1 {
		path = path[:len(path)-1]
	}
	return path
}

//Depth 从请求中读取depth头，并将其作为int类型返回
func Depth(r *http.Request) int {
	depth := Infinity
	depthStr := r.Header.Get("Depth")
	if len(depthStr) > 0 {
		if depthStr == "0" {
			depth = 0
		} else if depthStr == "1" {
			depth = 1
		}
	}
	return depth
}

func RewriteUrl(path string) string {
	return url.PathEscape(path)
}

//ETag 获取与文件关联的ETag。
func ETag(so *model.StoredObject) string {

	resourceLength := ""
	lastModified := ""
	if so != nil && !so.IsFolder {
		resourceLength = strconv.FormatInt(so.ContentLength, 10)
		lastModified = strconv.FormatInt(so.LastModified.Unix(), 10)
	}
	return "W/\"" + resourceLength + "-" + lastModified + "\""
}

func LockIdFromIfHeader(req *http.Request) []string {
	var ids []string
	id := req.Header.Get("If")

	if len(id) > 0 {
		if strings.Index(id, ">)") == strings.LastIndex(id, ">)") {
			id = id[strings.Index(id, "(<"):strings.Index(id, ">)")]
			if strings.Index(id, "locktoken:") != -1 {
				id = id[strings.Index(id, ":")+1:]
			}
			ids[0] = id
		} else {
			firstId := id[strings.Index(id, "(<"):strings.Index(id, (">)"))]
			if strings.Index(firstId, "locktoken:") != -1 {
				firstId = firstId[(strings.Index(firstId, ":") + 1):]
			}
			ids[0] = firstId

			secondId := id[strings.Index(id, "(<"):strings.Index(id, ">)")]
			if strings.Index(secondId, "locktoken:") != -1 {
				secondId = secondId[strings.Index(secondId, ":")+1:]
			}
			ids[1] = secondId
		}
	} else {
		ids = nil
	}
	return ids
}
func LockIdFromLockTokenHeader(req *http.Request) string {
	id := req.Header.Get("Lock-Token")
	if len(id) > 0 {
		id = id[strings.Index(id, ":")+1 : strings.Index(id, ">")]
	}
	return id
}

func CheckLocks(transaction model.Transaction, r *http.Request, w http.ResponseWriter
) bool {

}