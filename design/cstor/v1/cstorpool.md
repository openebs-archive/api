# Proposed v1 schema for CSPC and CSPI

## Table of Contents

* [Introduction](#introduction)
* [Goals](#goals)
* [CSPC and CSPI Schema Proposal](#cspc-and-cspi-schema-proposal)

## Introduction

This proposal highlights the limitations of current CSPC and CSPI schema and
proposes improvements for the same.
Please refer to the following link for cstor operator design document. 
https://github.com/openebs/openebs/blob/master/contribute/design/1.x/cstor-operator/doc.md

# Goals

OpenEBS cStor Data Engine has been in production for over a year now, with multiple cStor custom resource schema 
revisions. The goal of this design is to incorporate all the feedback received on the alpha and beta versions of cStor 
custom resources and define a version 1 schema.
The schema proposed in this document encompasses the following requirements:

- Include the best practicies followed by users with SPC into the schema itself.
- Rename the SPC and CSP to CSPC and CSPI based on the feedback received.
- Enhance the elements of CSPC and CSPI to enable performing Day 2 operations declaratively:
    - Scaling up/down of cStor Pools
    - Expanding the capacity of cStor Pools
    - Replacing the failed block devices in cStor Pools.
    - Add supportability features via Status Conditions and Events
    - Capture the history / log of upgrades


## CSPC and CSPI Schema Proposal

### CSPC Schema

```go
// CStorPoolCluster describes a CStorPoolCluster custom resource.
type CStorPoolCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CStorPoolClusterSpec   `json:"spec"`
	Status            CStorPoolClusterStatus `json:"status"`
	VersionDetails    VersionDetails         `json:"versionDetails"`
}

// CStorPoolClusterSpec is the spec for a CStorPoolClusterSpec resource
type CStorPoolClusterSpec struct {
	// Pools is the spec for pools for various nodes
	// where it should be created.
	Pools []PoolSpec `json:"pools"`

	// DefaultResources are the compute resources required by the cstor-pool
	// container.
	// If the resources at PoolConfig is not specified, this is written
	// to CSPI PoolConfig.
	DefaultResources *corev1.ResourceRequirements `json:"resources"`

	// AuxResources are the compute resources required by the cstor-pool pod
	// side car containers.
	DefaultAuxResources *corev1.ResourceRequirements `json:"auxResources"`

	// DefaultTolerations, if specified, are the pool pod's tolerations
	// If tolerations at PoolConfig is empty, this is written to
	// CSPI PoolConfig.
	DefaultTolerations []corev1.Toleration `json:"tolerations"`

	// DefaultPriorityClassName if specified applies to all the pool pods
	// in the pool spec if the priorityClass at the pool level is
	// not specified.
	DefaultPriorityClassName string `json:"priorityClassName"`
}

//PoolSpec is the spec for pool on node where it should be created.
type PoolSpec struct {
	// NodeSelector is the labels that will be used to select
	// a node for pool provisioning.
	// Required field
	NodeSelector map[string]string `json:"nodeSelector"`

	// DataRaidConfig is the raid group configuration for the given pool.
	DataRaidGroups []RaidGroup `json:"dataRaidGroups"`

	// WriteCacheGroups is the write cache for the given pool.
	WriteCacheGroups []RaidGroup `json:"writeCacheGroups"`

	// PoolConfig is the pool config that applies to the
	// pool on node.
	PoolConfig PoolConfig `json:"poolConfig"`
}

// PoolConfig is the default pool config that applies to the
// pool on node.
type PoolConfig struct {
	// DataRaidGroupType is the data raid group raid type 
	// Supported values are : stripe, mirror, raidz and raidz2
	DataRaidGroupType string `json:"dataRaidGroupType"`

	// WriteCacheRaidGroupType is the write cache raid type 
	// Supported values are : stripe, mirror, raidz and raidz2
	WriteCacheRaidGroupType string `json:"writeCacheRaidGroupType"`

	// ThickProvisioning will provision the volumes, only if capacity is available on the pool.
	// Setting this to true will disable over-provisioning of the pool.
	// Optional -- defaults to false
	ThickProvisioning bool `json:"thickProvisioning"`

	// Compression to enable compression
	// Optional -- defaults to off
	// Possible values : lz, off
	Compression string `json:"compression"`

	// Resources are the compute resources required by the cstor-pool
	// container.
	Resources *corev1.ResourceRequirements `json:"resources"`

	// AuxResources are the compute resources required by the cstor-pool pod
	// side car containers.
	AuxResources *corev1.ResourceRequirements `json:"auxResources"`

	// Tolerations, if specified, the pool pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations"`

	// PriorityClassName if specified applies to this pool pod
	// If left empty, DefaultPriorityClassName is applied.
	// (See CStorPoolClusterSpec.DefaultPriorityClassName)
	// If both are empty, not priority class is applied.
	PriorityClassName string `json:"priorityClassName"`
	
	// ROThresholdLimit is threshold(percentage base) limit
	// for pool read only mode. If ROThresholdLimit(%) amount
	// of pool storage is reached then pool will set to readonly.
	// NOTE:
	// 1. If ROThresholdLimit is set to 100 then entire
	//    pool storage will be used by default it will be set to 85%.
	// 2. ROThresholdLimit value will be 0 < ROThresholdLimit <= 100.
	// optional
	ROThresholdLimit int `json:"roThresholdLimit"`
}

// RaidGroup contains the details of a raid group for the pool
type RaidGroup struct {
	CStorPoolInstanceBlockDevices []CStorPoolInstanceBlockDevice `json:"cspiBlockDevices"`
}

// CStorPoolInstanceBlockDevice contains the details of block devices that
// constitutes a raid group.
type CStorPoolInstanceBlockDevice struct {
	// BlockDeviceName is the name of the block device.
	BlockDeviceName string `json:"blockDeviceName"`

	// Capacity is the capacity of the block device.
	// It is system generated filled by CSPI controller
	Capacity uint64 `json:"capacity"`

	// DevLink is the dev link for block devices
	DevLink string `json:"devLink"`
}

// CStorPoolClusterStatus represents the latest available observations of a CSPC's current state.
type CStorPoolClusterStatus struct {
	// ProvisionedInstances is the the number of CSPI present at the current state. 
	ProvisionedInstances int32 `json:"provisionedInstances"`
	
	// DesiredInstances is the number of CSPI(s) that should be provisioned.
	DesiredInstances int32 `json:"runningInstances"`

	// HealthyInstances is the number of CSPI(s) that are healthy.
	HealthyInstances int32 `json:"healthyInstances"`

	// Current state of CSPC.
	Conditions []CStorPoolClusterCondition  `json:conditions`
}

type CSPCConditionType string

// CStorPoolClusterCondition describes the state of a CSPC at a certain point.
type CStorPoolClusterCondition struct {
	// Type of CSPC condition.
	Type CSPCConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

// CStorPoolClusterList is a list of CStorPoolCluster resources
type CStorPoolClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CStorPoolCluster `json:"items"`
}

```

### CSPI Schema

```go
// CStorPoolInstance describes a cstor pool instance resource.
type CStorPoolInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Spec is the specification of the cstorpoolinstance resource.
	Spec CStorPoolInstanceSpec `json:"spec"`
	// Status is the possible statuses of the cstorpoolinstance resource.
	Status CStorPoolInstanceStatus `json:"status"`
	// VersionDetails is the openebs version.
	VersionDetails VersionDetails `json:"versionDetails"`
}

// CStorPoolInstanceSpec is the spec listing fields for a CStorPoolInstance resource.
type CStorPoolInstanceSpec struct {
	// HostName is the name of kubernetes node where the pool
	// should be created.
	HostName string `json:"hostName"`
	// NodeSelector is the labels that will be used to select
	// a node for pool provisioning.
	// Required field
	NodeSelector map[string]string `json:"nodeSelector"`
	// PoolConfig is the default pool config that applies to the
	// pool on node.
	PoolConfig PoolConfig `json:"poolConfig"`
	// DataRaidGroups is the raid group configuration for the given pool.
	DataRaidGroups []RaidGroup `json:"dataRaidGroups"`
	// WriteCacheRaidGroups is the write cache raid group.
	WriteCacheRaidGroups []RaidGroup `json:"writeCacheRaidGroups"`
}

// CStorPoolInstanceStatus is for handling status of pool.
type CStorPoolInstanceStatus struct {
	// Current state of CSPI with details.
	Conditions []CStorPoolInstanceCondition 
	// The phase of a CStorPool is a simple, high-level summary of the pool state on the node.  
	Phase    CStorPoolPhase        `json:"phase"`
	// Capacity describes the capacity details of a cstor pool 
	Capacity CStorPoolInstanceCapacity `json:"capacity"`
	//ReadOnly if pool is readOnly or not
	ReadOnly bool `json:"readOnly"`
	// ProvisionedReplicas describes the total count of Volume Replicas
	// present in the cstor pool
	ProvisionedReplicas int32 `json:"provisionedReplicas"`
	// HealthyReplicas describes the total count of healthy Volume Replicas
	// in the cstor pool
	HealthyReplicas int32 `json:"healthyReplicas"`
}

// CStorPoolInstanceCapacity stores the pool capacity related attributes.
type CStorPoolInstanceCapacity struct {
	// Amount of physical data (and its metadata) written to pool
	// after applying compression, etc..,
	Used resource.Quantity `json:"used"`
	// Amount of usable space in the pool after excluding
	// metadata and raid parity
	Free resource.Quantity `json:"free"`
	// Sum of usable capacity in all the data raidgroups
	Total resource.Quantity `json:"total"`
	// ZFSCapacityAttributes contains advanced information about pool capacity details
	ZFS ZFSCapacityAttributes `json:"zfs"`
}

// ZFSCapacityAttributes stores the advanced information about pool capacity related
// attributes
type ZFSCapacityAttributes struct {
	// LogicalUsed is the amount of space that is "logically" consumed
	// by this pool and all its descendents. The logical space ignores
	// the effect of the compression and copies properties, giving a
	// quantity closer to the amount of data that applications see.
	// However, it does include space consumed by metadata.
	// NOTE: Used and LogicalUsed can vary depends on the pool properties
	LogicalUsed resource.Quantity `json:"logicalUsed"`
}

type CSPIConditionType string

// CSPIConditionType describes the state of a CSPI at a certain point.
type CStorPoolInstanceCondition struct {
	// Type of CSPI condition.
	Type CSPIConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

```

