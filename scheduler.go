package operation_tools

import (
    "time"
)

type Scheduler struct {
    //计划任务列表
    ScheduleList []Schedule
}

type ScheduleInfo struct {
    Name string
    Desc string
}

//任务原型
type Schedule interface {
    //任务描述
    GetInfo() ScheduleInfo
    //运行任务
    Run(time.Time) error
}

//注册业务
func (p *Scheduler) Register(s Schedule) {
    p.ScheduleList = append(p.ScheduleList, s)
}

//获取当前计划任务列表
func (p *Scheduler) GetScheduleList() (arr []ScheduleInfo) {
    if 0 == len(p.ScheduleList) {
        return
    }

    arr = make([]ScheduleInfo, 0)
    for _, v := range p.ScheduleList {
        arr = append(arr, v.GetInfo())
    }
    return
}

//心跳
func (p *Scheduler) Heartbeat(nowTime time.Time) error {
    for _, v := range p.ScheduleList {
        err := v.Run(nowTime)
        if err != nil {
            return err
        }
    }
    return nil
}
