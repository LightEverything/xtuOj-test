package execode

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"xtuOj/models"
)

type ExeResult struct {
	Msg     string
	Total   int
	PassNum int
	Status  int
}

type ResultChan struct {
	Mem   chan int
	Time  chan int
	Pass  chan int
	CE    chan int
	WA    chan int
	MsgCE chan string
}

func ExeCodeWithTestCase(path string, input string, output string, rc *ResultChan, index int, maxMem int, complete chan bool) (e error) {
	// 初始化cmd的输入输出
	cmd := exec.Command("go", "run", path)
	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	io.WriteString(stdin, input)
	// 开始阶段的内存
	var beginMem runtime.MemStats
	runtime.ReadMemStats(&beginMem)

	if err := cmd.Run(); err != nil {
		if err.Error() == "exit status 1" {
			complete <- true
			rc.CE <- index
			rc.MsgCE <- stderr.String()
			return
		}
	}

	var endMem runtime.MemStats
	runtime.ReadMemStats(&endMem)

	if endMem.Alloc/1024-beginMem.Alloc/1024 > uint64(maxMem) {
		rc.Mem <- index
		complete <- true
		return
	}

	if out.String() == output {
		rc.Pass <- index
	} else {
		rc.WA <- index
	}

	// 通知检测程序已经完成
	complete <- true

	return nil
}

func ExeTestCase(ctx context.Context, rc *ResultChan, input string, output string, index int, maxMem int, path string) {
	complete := make(chan bool)

	go ExeCodeWithTestCase(path, input, output, rc, index, maxMem, complete)

	select {
	case <-ctx.Done():
		rc.Time <- index
		return
	case <-complete:
		return
	}
}

func ExeCode(ctx context.Context, data []models.Problem, path string) (e ExeResult) {
	// 获取问题
	pb := data[0]
	timeoutctx, cancel := context.WithTimeout(ctx, time.Duration(pb.MaxRuntime)*time.Millisecond)

	rc := ResultChan{
		Mem:   make(chan int),
		Time:  make(chan int),
		Pass:  make(chan int),
		CE:    make(chan int),
		WA:    make(chan int),
		MsgCE: make(chan string, 2),
	}

	// 执行每一个测试样例
	for i, v := range pb.TestCases {
		go ExeTestCase(timeoutctx, &rc, v.Input, v.Output, i, pb.MaxMem, path)
	}

	// 对每个样例单独判断
	er := ExeResult{}
	er.Total = len(pb.TestCases)
	caseStatus := make([]int, er.Total)

	for i := 0; i < len(pb.TestCases); i++ {
		select {
		case <-rc.CE:
			er.Status = 2
			er.PassNum = 0
			er.Msg = "编译错误！\n" + <-rc.MsgCE
			cancel()
			return er

		case i := <-rc.Pass:
			caseStatus[i] = 1
			er.PassNum++
		case i := <-rc.WA:
			caseStatus[i] = 2
		case i := <-rc.Time:
			caseStatus[i] = 3
		case i := <-rc.Mem:
			caseStatus[i] = 4
		}
	}

	// 构建msg返回值
	var msg strings.Builder

	for i, v := range caseStatus {
		switch v {
		case 1:
			msg.WriteString(fmt.Sprintf("测试样例%v通过\n", i))
		case 2:
			msg.WriteString(fmt.Sprintf("测试样例%v答案错误\n", i))
		case 3:
			msg.WriteString(fmt.Sprintf("测试样例%v超时\n", i))
		case 4:
			msg.WriteString(fmt.Sprintf("测试样例%v超出规定内存\n", i))
		}

	}

	if er.Total == er.PassNum {
		er.Status = 1
		er.Msg = "恭喜通过题目！\n" + msg.String()
	} else {
		er.Status = 2
		er.Msg = "题目未通过\n" + msg.String()
	}
	cancel()
	return er
}
