package master

import (
	"context"
	"crontab/master/common"
	"encoding/json"
	"time"

	v3 "go.etcd.io/etcd/client/v3"
)

// 任务管理
type JobMgr struct {
	client *v3.Client
	kv     v3.KV
	lease  v3.Lease
}

var (
	G_jobMgr *JobMgr
)

func InitJobMgr() (err error) {

	var (
		config v3.Config
		kv     v3.KV
		lease  v3.Lease
		client *v3.Client
	)

	config = v3.Config{
		Endpoints:   G_conf.EtcdEndpoints,
		DialTimeout: time.Duration(G_conf.EtcdDialTimeout) * time.Millisecond,
	}

	if client, err = v3.New(config); err != nil {
		return
	}

	kv = v3.NewKV(client)
	lease = v3.NewLease(client)

	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	return
}

func (JobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {

	var (
		jobKey    string
		jobValue  []byte
		putResp   *v3.PutResponse
		oldJobObj common.Job
	)

	jobKey = "/cron/jobs" + job.Name

	if jobValue, err = json.Marshal(job); err != nil {
		return
	}

	if putResp, err = JobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), v3.WithPrevKV()); err != nil {
		return
	}

	if putResp.PrevKv != nil {
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}
