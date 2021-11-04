package operation_tools

import (
    "bytes"
    "errors"
    "fmt"
    "os/exec"
    "strings"
)

func RunCmd(workDir string, cmdName string, args []string) (res string, err error) {
    cmdStr := fmt.Sprintf("%s %s", cmdName, strings.Join(args, " "))
    cmd := exec.Command("/bin/sh", "-c", cmdStr)
    if "" != workDir {
        cmd.Dir = workDir
    }
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    err = cmd.Run()
    if err != nil {
        err = errors.New(stderr.String())
        return
    }
    res = out.String()
    return
}
