// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package fileutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindFile(t *testing.T) {
	t.Run("files from various paths", func(t *testing.T) {
		// Create the following directory structure:
		// tmpDir1/
		//   file1.json
		//   file2.xml
		//   other.txt
		//   tmpDir2/
		//     other.txt/ [directory]
		//     tmpDir3/
		//       tmpDir4/
		//         tmpDir5/
		tmpDir1, err := ioutil.TempDir("", "")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir1)

		tmpDir2, err := ioutil.TempDir(tmpDir1, "")
		require.NoError(t, err)

		err = os.Mkdir(filepath.Join(tmpDir2, "other.txt"), 0700)
		require.NoError(t, err)

		tmpDir3, err := ioutil.TempDir(tmpDir2, "")
		require.NoError(t, err)

		tmpDir4, err := ioutil.TempDir(tmpDir3, "")
		require.NoError(t, err)

		tmpDir5, err := ioutil.TempDir(tmpDir4, "")
		require.NoError(t, err)

		type testCase struct {
			Description string
			Cwd         *string
			FileName    string
			Expected    string
		}

		testCases := []testCase{}

		for _, fileName := range []string{"file1.json", "file2.xml", "other.txt"} {
			filePath := filepath.Join(tmpDir1, fileName)
			require.NoError(t, ioutil.WriteFile(filePath, []byte("{}"), 0600))

			// Relative paths end up getting symlinks fully resolved, so use this below as necessary.
			filePathResolved, err := filepath.EvalSymlinks(filePath)
			require.NoError(t, err)

			testCases = append(testCases, []testCase{
				{
					fmt.Sprintf("absolute path to %s", fileName),
					nil,
					filePath,
					filePath,
				},
				{
					fmt.Sprintf("absolute path to %s from containing directory", fileName),
					&tmpDir1,
					filePath,
					filePath,
				},
				{
					fmt.Sprintf("relative path to %s from containing directory", fileName),
					&tmpDir1,
					fileName,
					filePathResolved,
				},
				{
					fmt.Sprintf("%s: subdirectory of containing directory", fileName),
					&tmpDir2,
					fileName,
					filePathResolved,
				},
				{
					fmt.Sprintf("%s: twice-nested subdirectory of containing directory", fileName),
					&tmpDir3,
					fileName,
					filePathResolved,
				},
				{
					fmt.Sprintf("%s: thrice-nested subdirectory of containing directory", fileName),
					&tmpDir4,
					fileName,
					filePathResolved,
				},
				{
					fmt.Sprintf("%s: can't find from four nesting levels deep", fileName),
					&tmpDir5,
					fileName,
					"",
				},
			}...)
		}

		for _, testCase := range testCases {
			t.Run(testCase.Description, func(t *testing.T) {
				if testCase.Cwd != nil {
					prevDir, err := os.Getwd()
					require.NoError(t, err)
					defer os.Chdir(prevDir)
					os.Chdir(*testCase.Cwd)
				}

				assert.Equal(t, testCase.Expected, FindFile(testCase.FileName))
			})
		}
	})
}
