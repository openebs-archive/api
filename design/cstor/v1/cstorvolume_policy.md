# Proposed v1 CStorVolumeConfig, CStorVolumePolicy and CStorVolumeAttachment Schema

## Table of Contents

* [Introduction](#introduction)
* [Goals](#goals)
* [VolumeConfig Schema Proposal](#volume-config-schema-proposal)
* [VolumePolicy Schema Proposal](#volume-policy-schema-proposal)
* [VolumeAttachment Schema Proposal](#volume-attachment-schema-proposal)

## Introduction

This proposal highlights enhancements and migration of current `CStorVolumeConfig` and `CStorVolumePolicy` schema and
proposes feature changes.

# Goals

The major goal of this document is to stabilize the APIs schema for `CStorVolumeConfig` and `CStorVolumePolicy`,
based on various learning and feedbacks for various volume related day 2 operations and how to make them easy to
perform.

# Volume Config Schema Proposal

## CStorVolumeConfig APIs

```go

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true

// CStorVolumeConfig describes a cstor volume config resource created as
// custom resource. CStorVolumeConfig is a request for creating cstor volume
// related resources like deployment, svc etc.
type CStorVolumeConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Spec defines a specification of a cstor volume config required
	// to provisione cstor volume resources
	Spec CStorVolumeConfigSpec `json:"spec"`

	// Publish contains info related to attachment of a volume to a node.
	// i.e. NodeId etc.
	Publish CStorVolumeConfigPublish `json:"publish,omitempty"`

	// Status represents the current information/status for the cstor volume
	// config, populated by the controller.
	Status         CStorVolumeConfigStatus `json:"status"`
	VersionDetails VersionDetails          `json:"versionDetails"`
}

// CStorVolumeConfigSpec is the spec for a CStorVolumeConfig resource
type CStorVolumeConfigSpec struct {
	// Capacity represents the actual resources of the underlying
	// cstor volume.
	Capacity corev1.ResourceList `json:"capacity"`
	// CStorVolumeRef has the information about where CstorVolumeClaim
	// is created from.
	CStorVolumeRef *corev1.ObjectReference `json:"cstorVolumeRef,omitempty"`
	// CStorVolumeSource contains the source volumeName@snapShotname
	// combaination.  This will be filled only if it is a clone creation.
	CStorVolumeSource string `json:"cstorVolumeSource,omitempty"`
	// Provision represents the initial volume configuration for the underlying
	// cstor volume based on the persistent volume request by user.
	Provision VolumeProvision `json:"provision"`
	// Policy contains volume specific required policies target and replicas
	Policy CStorVolumePolicySpec `json:"policy"`
}

type VolumeProvision struct {
	// Capacity represents initial capacity of volume replica required during
	// volume clone operations to maintain some metadata info related to child
	// resources like snapshot, cloned volumes.
	Capacity corev1.ResourceList `json:"capacity"`
	// ReplicaCount represents initial cstor volume replica count, its will not
	// be updated later on based on scale up/down operations, only readonly
	// operations and validations.
	ReplicaCount int `json:"replicaCount"`
}

// CStorVolumeConfigPublish contains info related to attachment of a volume to a node.
// i.e. NodeId etc.
type CStorVolumeConfigPublish struct {
	// NodeID contains publish info related to attachment of a volume to a node.
	NodeID string `json:"nodeId,omitempty"`
}

// CStorVolumeConfigPhase represents the current phase of CStorVolumeConfig.
type CStorVolumeConfigPhase string

const (
	//CStorVolumeConfigPhasePending indicates that the cvc is still waiting for
	//the cstorvolume to be created and bound
	CStorVolumeConfigPhasePending CStorVolumeConfigPhase = "Pending"

	//CStorVolumeConfigPhaseBound indiacates that the cstorvolume has been
	//provisioned and bound to the cstor volume config
	CStorVolumeConfigPhaseBound CStorVolumeConfigPhase = "Bound"

	//CStorVolumeConfigPhaseFailed indiacates that the cstorvolume provisioning
	//has failed
	CStorVolumeConfigPhaseFailed CStorVolumeConfigPhase = "Failed"
)

// CStorVolumeConfigStatus is for handling status of CstorVolume Claim.
// defines the observed state of CStorVolumeConfig
type CStorVolumeConfigStatus struct {
	// Phase represents the current phase of CStorVolumeConfig.
	Phase CStorVolumeConfigPhase `json:"phase"`

	// PoolInfo represents current pool names where volume replicas exists
	PoolInfo []string `json:"poolInfo"`

	// Capacity the actual resources of the underlying volume.
	Capacity corev1.ResourceList `json:"capacity,omitempty"`

	Conditions []CStorVolumeConfigCondition `json:"condition,omitempty"`
}

// CStorVolumeConfigCondition contains details about state of cstor volume
type CStorVolumeConfigCondition struct {
	// Current Condition of cstor volume config. If underlying persistent volume is being
	// resized then the Condition will be set to 'ResizeStarted' etc
	Type CStorVolumeConfigConditionType `json:"type"`
	// Last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason is a brief CamelCase string that describes any failure
	Reason string `json:"reason"`
	// Human-readable message indicating details about last transition.
	Message string `json:"message"`
}

// CStorVolumeConfigConditionType is a valid value of CstorVolumeConfigCondition.Type
type CStorVolumeConfigConditionType string

// These constants are CVC condition types related to resize operation.
const (
	// CStorVolumeConfigResizePending ...
	CStorVolumeConfigResizing CStorVolumeConfigConditionType = "Resizing"
	// CStorVolumeConfigResizeFailed ...
	CStorVolumeConfigResizeFailed CStorVolumeConfigConditionType = "VolumeResizeFailed"
	// CStorVolumeConfigResizeSuccess ...
	CStorVolumeConfigResizeSuccess CStorVolumeConfigConditionType = "VolumeResizeSuccessful"
	// CStorVolumeConfigResizePending ...
	CStorVolumeConfigResizePending CStorVolumeConfigConditionType = "VolumeResizePending"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true

// CStorVolumeConfigList is a list of CStorVolumeConfig resources
type CStorVolumeConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CStorVolumeConfig `json:"items"`
}
```


# Volume Policy Schema Proposal

## CStorVolumePolicy APIs

```go

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true

// CStorVolumePolicy describes a configuration required for cstor volume
// resources
type CStorVolumePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Spec defines a configuration info of a cstor volume required
	// to provisione cstor volume resources
	Spec   CStorVolumePolicySpec   `json:"spec"`
	Status CStorVolumePolicyStatus `json:"status"`
}

// CStorVolumePolicySpec ...
type CStorVolumePolicySpec struct {
	// replicaAffinity is set to true then volume replica resources need to be
	// distributed across the pool instances
	Provision Provision `json:"provision"`

	// TargetSpec represents configuration related to cstor target and its resources
	Target TargetSpec `json:"target"`

	// ReplicaSpec represents configuration related to replicas resources
	Replica ReplicaSpec `json:"replica"`

	// ReplicaPoolInfo holds the pool information of volume replicas.
	// Ex: If volume is provisioned on which CStor pool volume replicas exist
	ReplicaPoolInfo []ReplicaPoolInfo `json:"replicaPoolInfo"`
}

// TargetSpec represents configuration related to cstor target and its resources
type TargetSpec struct {
	// QueueDepth sets the queue size at iSCSI target which limits the
	// ongoing IO count from client
	QueueDepth string `json:"queueDepth,omitempty"`

	// IOWorkers sets the number of threads that are working on above queue
	IOWorkers int64 `json:"luWorkers,omitempty"`

	// Monitor enables or disables the target exporter sidecar
	Monitor bool `json:"monitor,omitempty"`

	// ReplicationFactor represents maximum number of replicas
	// that are allowed to connect to the target
	ReplicationFactor int64 `json:"replicationFactor,omitempty"`

	// Resources are the compute resources required by the cstor-target
	// container.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// AuxResources are the compute resources required by the cstor-target pod
	// side car containers.
	AuxResources *corev1.ResourceRequirements `json:"auxResources,omitempty"`

	// Tolerations, if specified, are the target pod's tolerations
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// PodAffinity if specified, are the target pod's affinities
	PodAffinity *corev1.PodAffinity `json:"affinity,omitempty"`

	// NodeSelector is the labels that will be used to select
	// a node for target pod scheduleing
	// Required field
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// PriorityClassName if specified applies to this target pod
	// If left empty, no priority class is applied.
	PriorityClassName string `json:"priorityClassName,omitempty"`
}

// ReplicaSpec represents configuration related to replicas resources
type ReplicaSpec struct {
	// IOWorkers represents number of threads that executes client IOs
	IOWorkers string `json:"zvolWorkers,omitempty"`
	// Controls the compression algorithm used for this volumes
	// examples: on|off|gzip|gzip-N|lz4|lzjb|zle
	//
	// Setting compression to "on" indicates that the current default compression
	// algorithm should be used.The default balances compression and decompression
	// speed, with compression ratio and is expected to work well on a wide variety
	// of workloads. Unlike all other set‐tings for this property, on does not
	// select a fixed compression type.  As new compression algorithms are added
	// to ZFS and enabled on a pool, the default compression algorithm may change.
	// The current default compression algorithm is either lzjb or, if the
	// `lz4_compress feature is enabled, lz4.

	// The lz4 compression algorithm is a high-performance replacement for the lzjb
	// algorithm. It features significantly faster compression and decompression,
	// as well as a moderately higher compression ratio than lzjb, but can only
	// be used on pools with the lz4_compress

	// feature set to enabled.  See zpool-features(5) for details on ZFS feature
	// flags and the lz4_compress feature.

	// The lzjb compression algorithm is optimized for performance while providing
	// decent data compression.

	// The gzip compression algorithm uses the same compression as the gzip(1)
	// command.  You can specify the gzip level by using the value gzip-N,
	// where N is an integer from 1 (fastest) to 9 (best compression ratio).
	// Currently, gzip is equivalent to gzip-6 (which is also the default for gzip(1)).

	// The zle compression algorithm compresses runs of zeros.
	Compression string `json:"compression,omitempty"`
}

// Provision represents different provisioning policy for cstor volumes
type Provision struct {
	// replicaAffinity is set to true then volume replica resources need to be
	// distributed across the cstor pool instances based on the given topology
	ReplicaAffinity bool `json:"replicaAffinity"`
	// BlockSize is the logical block size in multiple of 512 bytes
	// BlockSize specifies the block size of the volume. The blocksize
	// cannot be changed once the volume has been written, so it should be
	// set at volume creation time. The default blocksize for volumes is 4 Kbytes.
	// Any power of 2 from 512 bytes to 128 Kbytes is valid.
	BlockSize uint32 `json:"blockSize"`
}

// ReplicaPoolInfo represents the pool information of volume replica
type ReplicaPoolInfo struct {
	// PoolName represents the pool name where volume replica exists
	PoolName string `json:"poolName"`
	// UID also can be added
}

// CStorVolumePolicyStatus is for handling status of CstorVolumePolicy
type CStorVolumePolicyStatus struct {
	Phase string `json:"phase"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true

// CStorVolumePolicyList is a list of CStorVolumePolicy resources
type CStorVolumePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CStorVolumePolicy `json:"items"`
}

```

# Volume Attachment Schema Proposal

## CStorVolumeAttachment APIs
```go

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=csivolume

// CStorVolumeAttachment represents a CSI based volume
type CStorVolumeAttachment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CStorVolumeAttachmentSpec   `json:"spec"`
	Status CStorVolumeAttachmentStatus `json:"status"`
}

// CStorVolumeAttachmentSpec is the spec for a CStorVolume resource
type CStorVolumeAttachmentSpec struct {
	// Volume specific info
	Volume VolumeInfo `json:"volume"`

	// ISCSIInfo specific to ISCSI protocol,
	// this is filled only if the volume type
	// is iSCSI
	ISCSI ISCSIInfo `json:"iscsi"`
}

// VolumeInfo contains the volume related info
// for all types of volumes in CStorVolumeAttachmentSpec
type VolumeInfo struct {
	// Name of the CSI volume
	Name string `json:"name"`

	// Capacity of the volume
	Capacity string `json:"capacity,omitempty"`

	// OwnerNodeID is the Node ID which
	// is also the owner of this Volume
	OwnerNodeID string `json:"ownerNodeID"`

	// FSType of a volume will specify the
	// format type - ext4(default), xfs of PV
	FSType string `json:"fsType,omitempty"`

	// AccessMode of a volume will hold the
	// access mode of the volume
	AccessModes []string `json:"accessModes,omitempty"`

	// AccessType of a volume will indicate if the volume will be used as a
	// block device or mounted on a path
	AccessType string `json:"accessType,omitempty"`

	// StagingPath of the volume will hold the
	// path on which the volume is mounted
	// on that node
	StagingTargetPath string `json:"stagingTargetPath,omitempty"`

	// TargetPath of the volume will hold the
	// path on which the volume is bind mounted
	// on that node
	TargetPath string `json:"targetPath,omitempty"`

	// ReadOnly specifies if the volume needs
	// to be mounted in ReadOnly mode
	ReadOnly bool `json:"readOnly,omitempty"`

	// MountOptions specifies the options with
	// which mount needs to be attempted
	MountOptions []string `json:"mountOptions,omitempty"`

	// Device Path specifies the device path
	// which is returned when the iSCSI
	// login is successful
	DevicePath string `json:"devicePath,omitempty"`
}

// ISCSIInfo has ISCSI protocol specific info,
// this can be used only if the volume type exposed
// by the vendor is iSCSI
type ISCSIInfo struct {
	// Iqn of this volume
	Iqn string `json:"iqn"`

	// TargetPortal holds the target portal
	// of this volume
	TargetPortal string `json:"targetPortal"`

	// IscsiInterface of this volume
	IscsiInterface string `json:"iscsiInterface"`

	// Lun specify the lun number 0, 1.. on
	// iSCSI Volume. (default: 0)
	Lun string `json:"lun"`
}

// CStorVolumeAttachmentStatus status represents the current mount status of the volume
type CStorVolumeAttachmentStatus string

// CStorVolumeAttachmentStatusMounting indicated that a mount operation has been triggered
// on the volume and is under progress
const (
	// CStorVolumeAttachmentStatusUninitialized indicates that no operation has been
	// performed on the volume yet on this node
	CStorVolumeAttachmentStatusUninitialized CStorVolumeAttachmentStatus = ""
	// CStorVolumeAttachmentStatusMountUnderProgress indicates that the volume is busy and
	// unavailable for use by other goroutines, an iSCSI login followed by mount
	// is under progress on this volume
	CStorVolumeAttachmentStatusMountUnderProgress CStorVolumeAttachmentStatus = "MountUnderProgress"
	// CStorVolumeAttachmentStatusMounteid indicated that the volume has been successfulled
	// mounted on the node
	CStorVolumeAttachmentStatusMounted CStorVolumeAttachmentStatus = "Mounted"
	// CStorVolumeAttachmentStatusUnMounted indicated that the volume has been successfuly
	// unmounted and logged out of the node
	CStorVolumeAttachmentStatusUnmounted CStorVolumeAttachmentStatus = "Unmounted"
	// CStorVolumeAttachmentStatusRaw indicates that the volume is being used in raw format
	// by the application, therefore CSI has only performed iSCSI login
	// operation on this volume and avoided filesystem creation and mount.
	CStorVolumeAttachmentStatusRaw CStorVolumeAttachmentStatus = "Raw"
	// CStorVolumeAttachmentStatusResizeInProgress indicates that the volume is being
	// resized
	CStorVolumeAttachmentStatusResizeInProgress CStorVolumeAttachmentStatus = "ResizeInProgress"
	// CStorVolumeAttachmentStatusMountFailed indicates that login and mount process from
	// the volume has bben started but failed kubernetes needs to retry sending
	// nodepublish
	CStorVolumeAttachmentStatusMountFailed CStorVolumeAttachmentStatus = "MountFailed"
	// CStorVolumeAttachmentStatusUnmountInProgress indicates that the volume is busy and
	// unavailable for use by other goroutines, an unmount operation on volume
	// is under progress
	CStorVolumeAttachmentStatusUnmountUnderProgress CStorVolumeAttachmentStatus = "UnmountUnderProgress"
	// CStorVolumeAttachmentStatusWaitingForCVCBound indicates that the volume components
	// are still being created
	CStorVolumeAttachmentStatusWaitingForCVCBound CStorVolumeAttachmentStatus = "WaitingForCVCBound"
	// CStorVolumeAttachmentStatusWaitingForVolumeToBeReady indicates that the replicas are
	// yet to connect to target
	CStorVolumeAttachmentStatusWaitingForVolumeToBeReady CStorVolumeAttachmentStatus = "WaitingForVolumeToBeReady"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=csivolumes

// CStorVolumeAttachmentList is a list of CStorVolumeAttachment resources
type CStorVolumeAttachmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CStorVolumeAttachment `json:"items"`
}


```
