/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package scripting_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal/hammer/scripting"
)

var _ = Describe("WriteTempFile", func() {
	It("returns an empty file if given no input lines", func() {
		path, err := WriteTempFile()
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(path)

		contents, err := os.ReadFile(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(contents)).To(BeEmpty())
	})

	It("returns file containing input lines", func() {
		path, err := WriteTempFile("line1", "line2", "line3")
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(path)

		contents, err := os.ReadFile(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(contents)).To(Equal("line1\nline2\nline3\n"))
	})
})
