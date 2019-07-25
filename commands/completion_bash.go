/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package commands

import (
	"fmt"
)

type CompletionCommand struct {
	Bash BashCompletionCommand `command:"bash" description:"bash completion script"`
}

type BashCompletionCommand struct {
}

func (c *BashCompletionCommand) Execute(args []string) error {
	fmt.Printf(`# This script allows hammer to do autocompletion via Bash.
# Add the following to your .bashrc file, making sure that the path matches your system:
# eval "$(/path/to/hammer completion bash)"

_complete_hammer() {
  args=("${COMP_WORDS[@]:1:$COMP_CWORD}") # Skip first arg
  local IFS=$'\n' # Split into lines
  COMPREPLY=($(GO_FLAGS_COMPLETION=1 ${COMP_WORDS[0]} "${args[@]}"))
  return 0
}

complete -F _complete_hammer hammer
`)
	return nil
}
