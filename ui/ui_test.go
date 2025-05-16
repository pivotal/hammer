/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package ui_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	. "github.com/pivotal/hammer/ui"
)

var _ = Describe("ui", func() {
	var ui UI

	BeforeEach(func() {
		ui = UI{
			Out: gbytes.NewBuffer(),
			Err: gbytes.NewBuffer(),
		}
	})

	Context("DisplayText", func() {
		It("prints text to the out buffer", func() {
			ui.DisplayText("Test text output")
			Expect(ui.Out).To(gbytes.Say("Test text output"))
		})
	})

	Context("DisplayError", func() {
		It("prints error text to the err buffer", func() {
			ui.DisplayError(fmt.Errorf("test error contents"))
			Expect(ui.Err).To(gbytes.Say("test error contents"))
		})
	})
})
