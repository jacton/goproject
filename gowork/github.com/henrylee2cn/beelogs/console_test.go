// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package beelogs

import (
	"testing"
)

// Try each log level in decreasing order of priority.
func testConsoleCalls(bl *BeeLogger) {
	bl.Emergency("emergency")
	bl.Alert("alert")
	bl.Critical("critical")
	bl.Error("error")
	bl.Warning("warning")
	bl.Notice("notice")
	bl.Informational("informational")
	bl.Debug("debug")
}

// Test console logging by visually comparing the lines being output with and
// without a log level specification.
func TestConsole(t *testing.T) {
	log1 := NewLogger(10000)
	log1.EnableFuncCallDepth(true)
	log1.SetLogger("console", nil)
	testConsoleCalls(log1)

	log2 := NewLogger(100)
	log2.SetLogger("console", map[string]interface{}{"level": 3})
	testConsoleCalls(log2)
}

func BenchmarkConsole(b *testing.B) {
	log := NewLogger(10000)
	log.EnableFuncCallDepth(true)
	log.SetLogger("console", nil)
	for i := 0; i < b.N; i++ {
		log.Debug("debug")
	}
}