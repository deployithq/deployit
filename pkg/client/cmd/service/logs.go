//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package service

import (
	"errors"
	"fmt"
	e "github.com/lastbackend/lastbackend/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/unloop/gopipe"
	"io"
	"strings"
)

type mapInfo map[string]serviceInfo
type serviceInfo struct {
	Pod       string
	Container string
}

type Writer struct {
	io.Writer
}

func (Writer) Write(p []byte) (int, error) {
	return fmt.Print(string(p))
}

func LogsServiceCmd(name string) {

	var (
		ctx           = context.Get()
		choice string = "0"
	)

	service, projectName, err := Inspect(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	var m = make(mapInfo)
	//var index int = 0

	fmt.Println("Contaner list:\n")

	//for _, pod := range service.Spec.PodList {
	//	for _, container := range pod.ContainerList {
	//		fmt.Printf("[%d] %s\n", index, container.Name)
	//
	//		m[strconv.Itoa(index)] = serviceInfo{
	//			Pod:       pod.Name,
	//			Container: container.Name,
	//		}
	//	}
	//	index++
	//}

	if len(m) > 1 {
		for {
			fmt.Print("\nEnter container number for watch log or ^C for Exit: ")
			fmt.Scan(&choice)
			choice = strings.ToLower(choice)

			if _, ok := m[choice]; ok {
				break
			}

			ctx.Log.Error("Number not correct!")
		}
	}

	reader, err := Logs(projectName, service.Name, m[choice].Pod, m[choice].Container)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	stream.New(Writer{}).Pipe(reader)
}

func Logs(project, name, pod, container string) (*io.ReadCloser, error) {

	var (
		err error
		ctx = context.Get()
		er  = new(e.Http)
	)

	_, res, err := ctx.HTTP.
		GET("/project/"+project+"/service/"+name+"/logs?pod="+pod+"&container="+container).
		AddHeader("Authorization", "Bearer "+ctx.Token).Do()
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if er.Code == 401 {
		return nil, e.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, errors.New(er.Message)
	}

	return &res.Body, nil
}
