# Moving CSPC and CSPI Schema to cstor group and version v1 from openebs group and version v1alpha1

## Table of Contents

* [Introduction](#introduction)
* [Goals](#goals)
* [Current State of Things](#current-state-of-things)
* [CSPC and CSPI Schema Proposal](#cspc-schema-proposal)
* [CRD Migration](#crd-migration)

## Introduction

This proposal highlights the limitations of current CSPC and CSPI schema and
proposes improvements for the same.
Please refer to the following link for cstor operator design document. 
https://github.com/openebs/openebs/blob/master/contribute/design/1.x/cstor-operator/doc.md

# Goals

The major goal of this document is to freeze the schema for CSPC and CSPI schema.
Apart from this the document focuses on following aspects too:

- Moving the CSPC and CSPI CRD to v1 api version and cstor api group.
- Enhance the API in order to enhance user experience and support new features.


## Current State of Things

### CSPC Schema

The CSPC API has the following capabilities : 
 
- CSPC API can be used to provision cStor pools on a single or multiple nodes. Pool configuration for each node can be specified in the CSPC. 

- Stripe, mirror, raidz1, and raidz2 are the only supported raid topologies.

- A stripe raid group can have any number of block devices but not less than 1.

- A mirror, raidz and raidz2 raid group can only have exactly 2, 3 and 6 block devices only.  

- CSPC has the capability to specify a cache file for faster imports.

- CSPC has the capability to specify for overprovisioning and compression.

- CSPC has the capability to specify a default raid group type for a pool spec at the node level. If the raid group spec does not have a type -- this default raid group type is used.


- Resource requirements can be passed for cstor-pool and side-car containers via CSPC and there is a defaulting mechanism too. Please refer to the following PRs to understand more on this:
https://github.com/openebs/maya/pull/1444
https://github.com/openebs/maya/pull/1567
NOTE: If resource and limit requirements are mentioned nowhere in the CSPC then there is no default value that gets applied.


- CSPC can be used to specify pod priority for pool pods and there is a defaulting mechanism too. Please refer to the following PR to understand more.
https://github.com/openebs/maya/pull/1566


- CSPC can be used to specify tolerations for pool pods and there is a defaulting mechanism too. Please refer to the following PR to understand more.
https://github.com/openebs/maya/pull/1549


- CSPC can be used to do block device replacement.

- CSPC can be used to do pool expansion.

- CSPC can be used to create a pool on a new brought up node. ( Horizontal scale )

Following is the current CSPC schema in go struct : 

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
  // Tolerations, if specified, are the pool pod's tolerations
  // If tolerations at PoolConfig is empty, this is written to
  // CSPI PoolConfig.
  Tolerations []corev1.Toleration `json:"tolerations"`

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
  // RaidConfig is the raid group configuration for the given pool.
  RaidGroups []RaidGroup `json:"raidGroups"`
  // PoolConfig is the default pool config that applies to the
  // pool on node.
  PoolConfig PoolConfig `json:"poolConfig"`
}

// PoolConfig is the default pool config that applies to the
// pool on node.
type PoolConfig struct {
  // Cachefile is used for faster pool imports
  // optional -- if not specified or left empty cache file is not
  // used.
  CacheFile string `json:"cacheFile"`
  // DefaultRaidGroupType is the default raid type which applies
  // to all the pools if raid type is not specified there
  // Compulsory field if any raidGroup is not given Type
  DefaultRaidGroupType string `json:"defaultRaidGroupType"`

  // OverProvisioning to enable over provisioning
  // Optional -- defaults to false
  OverProvisioning bool `json:"overProvisioning"`
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
  
}

// RaidGroup contains the details of a raid group for the pool
type RaidGroup struct {
  // Type is the raid group type
  // Supported values are : stripe, mirror, raidz and raidz2

  // stripe -- stripe is a raid group which divides data into blocks and
  // spreads the data blocks across multiple block devices.

  // mirror -- mirror is a raid group which does redundancy
  // across multiple block devices.

  // raidz -- RAID-Z is a data/parity distribution scheme like RAID-5, but uses dynamic stripe width.
  // radiz2 -- TODO
  // Optional -- defaults to `defaultRaidGroupType` present in `PoolConfig`
  Type string `json:"type"`
  // IsWriteCache is to enable this group as a write cache.
  IsWriteCache bool `json:"isWriteCache"`
  // IsSpare is to declare this group as spare which will be
  // part of the pool that can be used if some block devices
  // fail.
  IsSpare bool `json:"isSpare"`
  // IsReadCache is to enable this group as read cache.
  IsReadCache bool `json:"isReadCache"`
  // BlockDevices contains a list of block devices that
  // constitute this raid group.
  BlockDevices []CStorPoolClusterBlockDevice `json:"blockDevices"`
}

// CStorPoolClusterBlockDevice contains the details of block devices that
// constitutes a raid group.
type CStorPoolClusterBlockDevice struct {
  // BlockDeviceName is the name of the block device.
  BlockDeviceName string `json:"blockDeviceName"`
  // Capacity is the capacity of the block device.
  // It is system generated
  Capacity string `json:"capacity"`
  // DevLink is the dev link for block devices
  DevLink string `json:"devLink"`
}

// CStorPoolClusterStatus is for handling status of pool.
type CStorPoolClusterStatus struct {
  Phase    string                `json:"phase"`
  Capacity CStorPoolCapacityAttr `json:"capacity"`
}

// CStorPoolClusterList is a list of CStorPoolCluster resources
type CStorPoolClusterList struct {
  metav1.TypeMeta `json:",inline"`
  metav1.ListMeta `json:"metadata"`

  Items []CStorPoolCluster `json:"items"`
}
```

### Limitations of Current CSPC

Although CSPC has an extensive schema, it suffers few limitations. I find following necessary improvements that can be done to make it more feature rich and pool related things more informative and debuggable.

Consider following improvement points : 

- Raidz can have 2+1, 4+1 or in general 2^n + 1 number of disks where n > 0.


- Similarly raidz2 can have 2^n + 2 number of disks where n>2. Although we can have 4 disks in raidz2 but that does not serve the purpose hence this should be restricted in the control plane.


- Write cache, read cache and spare property of a raid group cannot be specified simultaneously and there exists three different fields in CSPC to declare that. These three fields can be merged to be one field.


- CSPC does not have any status specifying the state of healthy pool instances vs total pool instances.


- There should exist a field to specify that the pool ( CSPI ) should not be considered for placing a cStor volume replica. ( Cordoning )


- The status field of CSPC should have following fields :
    - **currentProvisionedInstances:** Indicates the number of CSPI(s) present in the system.
	- **runningInstances:** Indicates the number of CSPI(s) that are not healthy but in degraded,rebuilding etc mode.
    - **healthyInstances:** Indicates the desired number of healthy CSPI(s) present in the system.
    - **conditions:** Represents the latest available observations of a CSPC’s current state. It is an array that can represent various conditions as an when required and hence will give more flexibility in terms of improving in areas of information and debuggability as and when required.


## CSPC and CSPI Schema Proposal

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
	// It is system generated
	Capacity string `json:"capacity"`

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

The CSPI resuses the CSPC schema and following is the current schema in go struct.

```go
// CStorPoolInstance describes a cstor pool instance resource created as custom resource.
type CStorPoolInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec           CStorPoolInstanceSpec `json:"spec"`
	Status         CStorPoolStatus       `json:"status"`
	VersionDetails VersionDetails        `json:"versionDetails"`
}

// CStorPoolInstanceSpec is the spec listing fields for a CStorPoolInstance resource.
type CStorPoolInstanceSpec struct {
	// HostName is the name of kubernetes node where the pool
	// should be created. This is filled by CSPC-operator.
	HostName string `json:"hostName"`
	// NodeSelector is the labels that is filled by CSPC-operator, only
	// for informational purpose.
	NodeSelector map[string]string `json:"nodeSelector"`
	// PoolConfig is the default pool config that applies to the
	// pool on node.
	PoolConfig PoolConfig `json:"poolConfig"`
	// RaidGroups is the group containing block devices
	RaidGroups []RaidGroup `json:"raidGroup"`
}

// CStorPoolStatus is for handling status of pool.
type CStorPoolStatus struct {
	Phase    CStorPoolPhase        `json:"phase"`
	Capacity CStorPoolCapacityAttr `json:"capacity"`
	// LastTransitionTime refers to the time when the phase changes
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	LastUpdateTime     metav1.Time `json:"lastUpdateTime,omitempty"`
	Message            string      `json:"message,omitempty"`
}

// CStorPoolCapacityAttr stores the pool capacity related attributes.
type CStorPoolCapacityAttr struct {
	Total string `json:"total"`
	Free  string `json:"free"`
	Used  string `json:"used"`
}

// CStorPoolInstanceList is a list of CStorPoolInstance resources
type CStorPoolInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CStorPoolInstance `json:"items"`
}

```

Following improvments can be made to the schema :
- The capacity fields ( i.e. Total, Free and Used) can be made to use Qunatity type that is helpful in capacity parsing and comparisons.
- Status can be improved to include `Conditions` which can describe the state of pool in more detail.

Following the the proposed schema for Status and Capacity (i.e. CStorPoolStatus and CStorPoolCapacityAttr struct ).
(Also renaming th struct CStorPoolStatus to CStorPoolInstanceStatus and CStorPoolCapacityAttr to CStorPoolInstanceCapacity)

```go

// CStorPoolInstanceStatus is for handling status of pool.
type CStorPoolInstanceStatus struct {
	// Current state of CSPI with details.
	Conditions []CStorPoolInstanceCondition 
	// The phase of a CStorPool is a simple, high-level summary of the pool state on the node.  
	Phase    CStorPoolPhase        `json:"phase"`
	// Capacity describes the capacity details of a cstor pool 
	Capacity CStorPoolInstanceCapacity `json:"capacity"`
}

// CStorPoolInstanceCapacity stores the pool capacity related attributes.
type CStorPoolInstanceCapacity struct {
	Total resource.Quantity `json:"total"`
	Free  resource.Quantity `json:"free"`
	Used  resource.Quantity `json:"used"`
}
type CSPIConditionType string

// CSPIConditionType describes the state of a CSPI at a certain point.
type CStorPoolInstanceCondition struct {
	// Type of CSPI condition.
	Type CSPIConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

```
Above schema has following advantages over the older one :
- Capacity parsing and comparison is easy.
- More details about the state of a pool.

