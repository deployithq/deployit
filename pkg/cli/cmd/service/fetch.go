//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/cli/view"
	"github.com/spf13/cobra"
)

func FetchCmd(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		cmd.Help()
		return
	}

	var namespace string
	cmd.Flags().StringVarP(&namespace, "namespace", "ns", "", "namespace")

	cli := context.Get().GetClient()
	response, err := cli.V1().Namespace(namespace).Service(args[0]).Get(context.Background())
	if err != nil {
		return
	}

	ss := view.FromApiServiceView(response)
	ss.Print()
}