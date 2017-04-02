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

package types

import "time"

type ContainerSpec struct {
	// Image spec
	Image         ImageSpec `json:"image"`
	// Network spec
	Network       ContainerNetworkSpec `json:"network"`
	// Ports configuration
	Ports         []ContainerPortSpec `json:"ports"`
	// Labels list
	Labels        map[string]string `json:"labels"`
	// Environments list
	Envs          map[string]string `json:"envs"`
	// Container enrtypoint
	Entrypoint    string `json:"entrypoint"`
	// Container run command
	Command       string `json:"command"`
	// Container run command arguments
	Args          []string `json:"args"`
	// Container DNS configuration
	DNS           ContainerDNSSpec `json:"dns"`
	// Container resources quota
	Quota         ContainerQuotaSpec `json:"quota"`
	// Container restart policy
	RestartPolicy ContainerRestartPolicySpec `json:"restart_policy"`
	// Container volumes mount
	Volumes       []VolumesSpec `json:"volumes"`
}


type ContainerNetworkSpec struct {
	// Container hostname
	Hostname string `json:"hostname"`
	// Container host domain
	Domain   string `json:"domain"`
	// Network ID to use
	Network  string `json:"network"`
}

type ContainerPortSpec struct {
	// Container port to expose
	ContainerPort int `json:"container_port"`
	// Containers protocol allowed on exposed port
	Protocol string `json:"protocol"`
}

type ContainerDNSSpec struct {
	// List of DNS servers
	Server  []string `json:"server"`
	// DNS server search options
	Search  []string `json:"search"`
	// DNS server other options
	Options []string `json:"options"`
}

type ContainerQuotaSpec struct {
	// Maximum memory allowed to use
	Memory int `json:"memory"`
	// CPU shares for container on one node
	CPUShares int `json:"cpu_shares"`
}

type ContainerRestartPolicySpec struct {
	// Restart policy name
	Name string `json:"name"`
	// Attempt to restart container
	Attempt int `json:"attempt"`
}

type ContainerStatusInfo struct {
	// Container ID on host
	ID string `json:"cid"`
	// Container current state
	State string `json:"state"`
	// Container ports mapping
	Ports map[string]int
	// Container created time
	Created time.Time `json:"created"`
	// Container updated time
	Updated time.Time `json:"updated"`
}