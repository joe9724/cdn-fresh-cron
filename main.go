package main

import (

_"github.com/jinzhu/gorm"
_ "github.com/jinzhu/gorm/dialects/mysql"
"fmt"
"cdn-fresh-cron/model"
"cdn-fresh-cron/utils"
"github.com/jasonlvhit/gocron"
	"time"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"unsafe"
	"encoding/json"
	"bytes"
	"net/http"
	"strings"
	"log"
)

func main() {
	//查询数据库，查找出需要刷新的url
	db, err := utils.OpenConnection()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	gocron.Every(10).Seconds().Do(getdata,db)


}

func getdata(db *gorm.DB){
	//先清空队列
	gocron.Clear()
	var produces []model.ProduceModel
	db.Raw("select * from produce where status=0 or status=1 or status=3").Find(&produces)  //3:刷新失败
	//判断上一次执行时间
	for k:=0; k<len(produces);k++  {
		//fmt.Println("try_times is:",produces[k].Status,produces[k].TryTimes)
		//根据try_times第几次,开启timer定时器
		//var delay int64
		if produces[k].TryTimes == 1{
			//delay = 10
			gocron.Every(1).Seconds().Do(task,db,produces[k].Status,produces[k].TryTimes,produces[k].TaskID)
		}else if produces[k].TryTimes == 2{
			//delay = 60
			gocron.Every(5).Seconds().Do(task,db,produces[k].Status,produces[k].TryTimes,produces[k].TaskID)
		}else if produces[k].TryTimes == 3{
			//delay = 120
			gocron.Every(10).Seconds().Do(task,db,produces[k].Status,produces[k].TryTimes,produces[k].TaskID)
		}
		//fmt.Println("delay is",delay)

	}
	<-gocron.Start()
}

func task(db *gorm.DB,status int64,trytimes int64,task_id string){
	fmt.Println("params is",status,trytimes)
	//对又拍发起轮询
	QueryYoupai(db,status,trytimes,task_id)
}

func QueryYoupai(db *gorm.DB, status int64,trytimes int64,task_id string) {
	host := "https://api.upyun.com/purge?task_ids="+task_id
	fmt.Println("get host is:",host)
	request, err := http.NewRequest("GET", host, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Authorization", "Bearer a542457c-cf65-40da-902a-a3c66100d063")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println(*str)
	fmt.Println("result is:", string(respBytes))

	//判断返回的code
	//如果有error_code
	if strings.Contains(string(respBytes), "error_code") {
		var errormodel model.ErrorResult
		json.Unmarshal(respBytes, &errormodel)
		//记录下error
		log.Println(string(respBytes))

	} else {
		//正常流程
		var result model.QueryModel
		json.Unmarshal(respBytes, &result)
		if len(result.Result)>0{
			if result.Result[0].Progress == 100 { //1代表刷新成功
				fmt.Println("update produce set status=2,last_update=?,progress=100 where task_id=?", time.Now().Unix(), task_id)
				db.Exec("update produce set status=2,last_update=?,progress=100 where task_id=?", time.Now().Unix(), task_id)
				//remove 定时任务
				//gocron.Remove()
			} else { //刷新中，进度不到100%
				fmt.Println("update produce set status=?,last_update=?,progress=? where task_id=?", 1, time.Now().Unix(), result.Result[0].Progress,task_id)
				db.Exec("update produce set status=?,last_update=?,progress=? where task_id=?", 1, time.Now().Unix(), result.Result[0].Progress,task_id)
			}
			//回调刷新结果给泰德
			var resultmodel model.ResultModel
			db.Raw("select * from produce where task_id=?",task_id).Find(&resultmodel)
			CallbackToTD(resultmodel.EventID, resultmodel.Status)
		}else{
			//回调刷新结果给泰德
			db.Exec("update produce set status=2 where task_id=?",task_id)
			var resultmodel model.ResultModel
			db.Raw("select * from produce where task_id=?",task_id).Find(&resultmodel)
			CallbackToTD(resultmodel.EventID, "2")
		}


	}





}
func CallbackToTD(event_id string, status string) {
	//send callback to td
	json_content := make(map[string]interface{})
	json_content["event_id"] = event_id
	json_content["status"] = status

	bytesData, err := json.Marshal(json_content)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	TDHost := "xxx"
	request, err := http.NewRequest("POST", TDHost, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	//request.Header.Set("Authorization","Bearer a542457c-cf65-40da-902a-a3c66100d0631")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("return from td is", string(respBytes))

	//get result from td

	//如果TD确认接收完毕 update database Mysql


}