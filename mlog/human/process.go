// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package human

import (
	"bufio"
	"io"
)

type LogWriter interface {
	Write(e LogEntry)
}

// Read JSON logs from input and write formatted logs to the output
func ProcessLogs(reader io.Reader, writer LogWriter) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		s := scanner.Text()
		e := ParseLogMessage(s)
		writer.Write(e)
	}
}
