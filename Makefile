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

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
GOBIN := $(or $(shell go env GOBIN 2>/dev/null), $(shell go env GOPATH 2>/dev/null)/bin)

# find or download controller-gen
controller-gen:
ifneq ($(shell controller-gen --version 2> /dev/null), Version: v0.11.4)
	@(cd /tmp;go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.11.4)
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

# Generate code, CRDs and documentation
ALL_CRDS=config/crds/all-crds.yaml
generate: generate-crds

generate-crds: controller-gen
	# Generate manifests e.g. CRD, RBAC etc.
	$(CONTROLLER_GEN) crd:crdVersions=v1 paths="./pkg/apis/cstor/..." output:crd:artifacts:config=config/crds/bases
	$(CONTROLLER_GEN) crd:crdVersions=v1 paths="./pkg/apis/openebs.io/v1alpha1/..." output:crd:artifacts:config=config/crds
	# merge all crds into a single file
	rm $(ALL_CRDS)
	cat config/crds/bases/*.yaml >> $(ALL_CRDS)

.PHONY: kubegen
# code generation for custom resources
kubegen:
	./hack/update-codegen.sh
	@cp -rf v3/pkg/apis pkg/
	@cp -rf v3/pkg/client pkg/
	@cp -rf v3/pkg/internalapis/apis/cstor/v1/zz_generated.conversion.go pkg/internalapis/apis/cstor/v1/zz_generated.conversion.go
	@cp -rf v3/pkg/internalapis/apis/cstor/zz_generated.deepcopy.go pkg/internalapis/apis/cstor/zz_generated.deepcopy.go
	@rm -rf v3

.PHONY: verify_kubegen
verify_kubegen:
	./hack/verify-codegen.sh

.PHONY: generated_files
generated_files: kubegen protobuf

.PHONY: protobuf
protobuf:
	@echo "----------------------------"
	@echo "--> protobuf           "
	@echo "----------------------------"
	@protoc -I $(PWD)/pkg/apis/openebs.io/v1alpha1/ \
    -I ${GOPATH}/src \
    --go_out=$(PWD)/pkg/proto \
    --go-grpc_out=require_unimplemented_servers=false:$(PWD)/pkg/proto \
    $(PWD)/pkg/apis/cstor/v1/cstorvolume.proto

# deps ensures fresh go.mod and go.sum.
.PHONY: deps
deps:
	@go mod tidy
	@go mod verify

.PHONY: test
test:
	go test ./...
	
.PHONY: license-check
license-check:
	@echo "--> Checking license header..."
	@licRes=$$(for file in $$(find . -type f -regex '.*\.sh\|.*\.go\|.*Docker.*\|.*\Makefile*') ; do \
               awk 'NR<=5' $$file | grep -Eq "(Copyright|generated|GENERATED)" || echo $$file; \
       done); \
       if [ -n "$${licRes}" ]; then \
               echo "license header checking failed:"; echo "$${licRes}"; \
               exit 1; \
       fi
	@echo "--> Done checking license."
	@echo
