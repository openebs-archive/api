/*
Copyright 2020 The OpenEBS Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import 	(
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/openebs/maya/pkg/util"

)

const (
	StoragePoolKindCSPC = "CStorPoolCluster"
	// APIVersion holds the value of OpenEBS version
	APIVersion = "cstor.openebs.io/v1"
)

func NewCStorPoolInstance()*CStorPoolInstance  {
	return &CStorPoolInstance{}
}


// WithName sets the Name field of CSPI with provided value.
func (cspi *CStorPoolInstance) WithName(name string) *CStorPoolInstance {
	cspi.Name = name
	return cspi
}

// WithNamespace sets the Namespace field of CSPI provided arguments
func (cspi *CStorPoolInstance) WithNamespace(namespace string) *CStorPoolInstance {
	cspi.Namespace = namespace
	return cspi
}

// WithAnnotationsNew sets the Annotations field of CSPI with provided arguments
func (cspi *CStorPoolInstance) WithAnnotationsNew(annotations map[string]string) *CStorPoolInstance {
	cspi.Annotations = make(map[string]string)
	for key, value := range annotations {
		cspi.Annotations[key] = value
	}
	return cspi
}

// WithAnnotations appends or overwrites existing Annotations
// values of CSPI with provided arguments
func (cspi *CStorPoolInstance) WithAnnotations(annotations map[string]string) *CStorPoolInstance {

	if cspi.Annotations == nil {
		return cspi.WithAnnotationsNew(annotations)
	}
	for key, value := range annotations {
		cspi.Annotations[key] = value
	}
	return cspi
}

// WithLabelsNew sets the Labels field of CSPI with provided arguments
func (cspi *CStorPoolInstance)WithLabelsNew(labels map[string]string) *CStorPoolInstance {
	cspi.Labels = make(map[string]string)
	for key, value := range labels {
		cspi.Labels[key] = value
	}
	return cspi
}

// WithLabels appends or overwrites existing Labels
// values of CSPI with provided arguments
func (cspi *CStorPoolInstance) WithLabels(labels map[string]string) *CStorPoolInstance {
	if cspi.Labels == nil {
		return cspi.WithLabelsNew(labels)
	}
	for key, value := range labels {
		cspi.Labels[key] = value
	}
	return cspi
}

// WithNodeSelectorByReference sets the node selector field of CSPI with provided argument.
func (cspi *CStorPoolInstance) WithNodeSelectorByReference(nodeSelector map[string]string) *CStorPoolInstance {
	cspi.Spec.NodeSelector = nodeSelector
	return cspi
}

// WithNodeName sets the HostName field of CSPI with the provided argument.
func (cspi *CStorPoolInstance) WithNodeName(nodeName string) *CStorPoolInstance {
	cspi.Spec.HostName = nodeName
	return cspi
}

// WithPoolConfig sets the pool config field of the CSPI with the provided config.
func (cspi *CStorPoolInstance) WithPoolConfig(poolConfig PoolConfig) *CStorPoolInstance {
	cspi.Spec.PoolConfig = poolConfig
	return cspi
}

// WithRaidGroups sets the raid group field of the CSPI with the provided raid groups.
func (cspi *CStorPoolInstance) WithRaidGroups(raidGroup []RaidGroup) *CStorPoolInstance {
	cspi.Spec.RaidGroups = raidGroup
	return cspi
}

// WithFinalizer sets the finalizer field in the BDC
func (cspi *CStorPoolInstance) WithFinalizer(finalizers ...string) *CStorPoolInstance {
	cspi.Finalizers = append(cspi.Finalizers, finalizers...)
	return cspi
}


// WithCSPCOwnerReference sets the OwnerReference field in CSPI with required
//fields
func (cspi *CStorPoolInstance) WithCSPCOwnerReference(reference metav1.OwnerReference) *CStorPoolInstance {
	cspi.OwnerReferences = append(cspi.OwnerReferences, reference)
	return cspi
}

// WithNewVersion sets the current and desired version field of
// CSPI with provided arguments
func (cspi *CStorPoolInstance) WithNewVersion(version string) *CStorPoolInstance {
	cspi.VersionDetails.Status.Current = version
	cspi.VersionDetails.Desired = version
	return cspi
}

// WithDependentsUpgraded sets the field to true for new CSPI
func (cspi *CStorPoolInstance) WithDependentsUpgraded() *CStorPoolInstance {
	cspi.VersionDetails.Status.DependentsUpgraded = true
	return cspi
}

// HasFinalizer returns true if the provided finalizer is present on the object.
func (cspi *CStorPoolInstance) HasFinalizer(finalizer string) bool {
	finalizersList := cspi.GetFinalizers()
	return util.ContainsString(finalizersList, finalizer)
}

// RemoveFinalizer removes the given finalizer from the object.
func (cspi *CStorPoolInstance) RemoveFinalizer(finalizer string)  {
	cspi.Finalizers = util.RemoveString(cspi.Finalizers, finalizer)
}

// HasAnnotation return true if provided annotation
// key and value are present on the object.
func (cspi *CStorPoolInstance) HasAnnotation(key, value string) bool {
	val, ok := cspi.GetAnnotations()[key]
	if ok {
		return val == value
	}
	return false
}

// HasLabel returns true if provided label
// key and value are present on the object.
func (cspi *CStorPoolInstance) HasLabel(key, value string) bool {
	val, ok := cspi.GetLabels()[key]
	if ok {
		return val == value
	}
	return false
}

