// Copyright © 2017 huang jia <449264675@qq.com>
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

package main

import (
	"flag"

	"server/cmd/apiserver/app"
	"server/pkg/utils/log"
)

var config = flag.String("config", "system.ini", "absolute config file path")

func main() {
	flag.Parse()
	s := app.NewAPIServer(*config)
	if err := app.Run(s); err != nil {
		log.Critical("start apiserver err: %v", err)
	}
}
