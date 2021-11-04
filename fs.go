package operation_tools

import (
    "crypto/md5"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
)

//获取文件列表
func GetFileList(rootDir string, ignoreMap map[string]int) (fileList []string, err error) {
    list, err := ioutil.ReadDir(rootDir)
    if err != nil {
        return
    }
    for _, v := range list {
        fileName := v.Name()
        abPath := fmt.Sprintf("%s/%s", rootDir, fileName)

        if isPathIgnored(abPath, ignoreMap) {
            continue
        }

        //如果是目录则递归
        if v.IsDir() {
            subFileList, err1 := GetFileList(abPath, ignoreMap)
            if err1 != nil {
                return
            }
            fileList = append(fileList, subFileList...)
        } else {
            //收集结果
            fileList = append(fileList, abPath)
        }
    }
    return
}

//检查代码dst目录是否存在
func IsDirExist(dir string) (isExist bool, err error) {
    _, err = os.Stat(dir)
    if err == nil {
        return true, nil
    }
    if os.IsExist(err) {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    //其他错误
    return false, err
}

//复制文件
func CopyDir(from string, to string) (err error) {
    from = fmt.Sprintf(`%s`, strings.TrimRight(from, "/")+"/*")
    to = strings.TrimRight(to, "/") + "/"
    _, err = RunCmd("", "/bin/cp", []string{"-r", from, to})
    return
}

//计算文件的md5
func CalcFileMd5(filePath string) (res string, err error) {
    body, err := ioutil.ReadFile(filePath)
    if err != nil {
        return
    }
    res = fmt.Sprintf("%x", md5.Sum(body))
    return
}

//路径是否被忽略
func isPathIgnored(abPath string, ignoreMap map[string]int) bool {
    for k := range ignoreMap {
        if isPathIgnoredByRule(abPath, k) {
            return true
        }
    }
    return false
}

func isPathIgnoredByRule(abPath string, rule string) bool {
    //以/开头 是绝对路径
    isAbsolute := "/" == rule[0:1]
    //以*结尾 是通配符
    isPattern := "*" == rule[len(rule)-1:]

    //通配符
    if isPattern {
        //去掉*
        ignorePath := rule[0 : len(rule)-1]
        if isAbsolute {
            //绝对路径通配符
            return 0 == strings.Index(abPath, ignorePath)
        } else {
            //任意路径通配符 先补齐左斜杠(防止出现无脑截取对比的情况)
            ignorePath = "/" + strings.TrimLeft(ignorePath, "/")
            return strings.Index(abPath, ignorePath) >= 0
        }
    }

    //非通配符
    if isAbsolute {
        //绝对路径通配符
        return abPath == rule
    } else {
        //任意路径通配符 先补齐左斜杠(防止出现无脑截取对比的情况) 然后不加右斜杠对齐尾部 或加上右斜杠能寻找到
        if strings.Index(abPath, "/"+rule+"/") >= 0 {
            return true
        } else {
            ignorePath := "/" + rule
            return len(ignorePath)+strings.Index(abPath, ignorePath) == len(abPath)
        }
    }
}
