# Copyright © 2020 The OpenEBS Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.



# Default arguments for code gen script

# OUTPUT_PKG is the path of directory where you want to keep the generated code
OUTPUT_PKG=github.com/openebs/api/pkg/client/generated

# APIS_PKG is the path where apis group and schema exists.
APIS_PKG=github.com/openebs/api/pkg/apis
INTERNAL_APIS_PKG=github.com/openebs/api/pkg/internal/apis

# GENS is an argument which generates different type of code.
# Possible values: all, deepcopy, client, informers, listers.
GENS=all
# GROUPS_WITH_VERSIONS is the group containing different versions of the resources.
GROUPS_WITH_VERSIONS="cstor:v1 openebs.io:v1alpha1"
INTERNAL_GROUPS_WITH_VERSIONS="cstor:v1"

# BOILERPLATE_TEXT_PATH is the boilerplate text(go comment) that is put at the top of every generated file.
# This boilerplate text is nothing but the license information.
BOILERPLATE_TEXT_PATH=hack/custom-boilerplate.go.txt

.PHONY: kubegen
# code generation for custom resources
kubegen: 
	./hack/code-gen.sh ${GENS} ${OUTPUT_PKG} ${APIS_PKG} ${GROUPS_WITH_VERSIONS} --go-header-file ${BOILERPLATE_TEXT_PATH}

.PHONY: generated_files
generated_files: kubegen protobuf

