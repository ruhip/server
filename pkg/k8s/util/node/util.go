// Copyright © 2018 huang jia <449264675@qq.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package node

import (
	"k8s.io/api/core/v1"
)

// IsNodeReady returns true if a node is ready; false otherwise.
func IsNodeReady(node *v1.Node) bool {
	for _, c := range node.Status.Conditions {
		if c.Type == v1.NodeReady {
			return c.Status == v1.ConditionTrue
		}
	}
	return false
}

// IsNodeNoSchedule returns true if a node is schedulable; false otherwise.
func IsNodeSchedule(node *v1.Node) bool {
	for _, t := range node.Spec.Taints {
		if t.Effect == v1.TaintEffectNoSchedule || node.Spec.Unschedulable {
			return false
		}
	}
	return true
}

// IsDiskPressure returns true if a node is DiskPressure; false otherwise.
func IsDiskPressure(node *v1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeDiskPressure && condition.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

// IsMemoryPressure returns true if a node is MemoryPressure; false otherwise.
func IsMemoryPressure(node *v1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeMemoryPressure && condition.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

func GetHostName(node *v1.Node) string {
	for _, address := range node.Status.Addresses {
		if address.Type == v1.NodeHostName {
			return address.Address
		}
	}
	return ""
}

func GetInternalIP(node *v1.Node) string {
	for _, addresse := range node.Status.Addresses {
		if addresse.Type == v1.NodeInternalIP {
			return addresse.Address
		}
	}
	return ""
}
