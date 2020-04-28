# Proposed v1 CStorVolumes and CStorVolumeReplica Schema

## Table of Contents

* [Introduction](#introduction)
* [Goals](#goals)
* [CStor Volume Schema Proposal](#cstor-volume-schema-proposal)

## Introduction

This proposal highlights the enhancements and migration of current `CStorVolumes` and `CStorVolumeReplicas` schema and
proposes improvements.

# Goals

The major goal of this document is to freeze the APIs schema for `CStorVolumes` and `CStorVolumeReplica`,
based on various learning and feedbacks from the user including various volume related day 2 operations.

# CStor Volume Schema Proposal

## `CStorVolumes` APIs

```go
// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CStorVolume describes a cstor volume resource created as custom resource
type CStorVolume struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CStorVolumeSpec   `json:"spec"`
	VersionDetails    VersionDetails    `json:"versionDetails"`
	Status            CStorVolumeStatus `json:"status"`
}


// CStorVolumeSpec is the spec for a CStorVolume resource
type CStorVolumeSpec struct {
	// Capacity represents the desired size of the underlying volume.
	Capacity resource.Quantity `json:"capacity"`

	// TargetIP IP of the iSCSI target service
	TargetIP string `json:"targetIP"`

	// iSCSI Target Port typically TCP ports 3260
	TargetPort string `json:"targetPort"`

	// Target iSCSI Qualified Name.combination of nodeBase
	Iqn string `json:"iqn"`

	// iSCSI Target Portal. The Portal is combination of IP:port (typically TCP ports 3260)
	TargetPortal string `json:"targetPortal"`

	// ReplicationFactor represents number of volume replica created during volume
	// provisioning connect to the target
	ReplicationFactor int `json:"replicationFactor"`

	// ConsistencyFactor is minimum number of volume replicas i.e. `RF/2 + 1`
	// has to be connected to the target for write operations. Basically more then
	// 50% of replica has to be connected to target.
	ConsistencyFactor int `json:"consistencyFactor"`

	// DesiredReplicationFactor represents maximum number of replicas
	// that are allowed to connect to the target. Required for scale operations
	DesiredReplicationFactor int `json:"desiredReplicationFactor"`

	//ReplicaDetails refers to the trusty replica information
	ReplicaDetails CStorVolumeReplicaDetails `json:"replicaDetails,omitempty"`
}

// ReplicaID is to hold replicaID information
type ReplicaID string

// CStorVolumePhase is to hold result of action.
type CStorVolumePhase string

// CStorVolumeStatus is for handling status of cvr.
type CStorVolumeStatus struct {
	Phase           CStorVolumePhase `json:"phase"`
	ReplicaStatuses []ReplicaStatus  `json:"replicaStatuses,omitempty"`
	// Represents the actual resources of the underlying volume.
	Capacity resource.Quantity `json:"capacity,omitempty"`
	// LastTransitionTime refers to the time when the phase changes
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// LastUpdateTime refers to the time when last status updated due to any
	// operations
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// A human-readable message indicating details about why the volume is in this state.
	Message string `json:"message,omitempty"`
	// Current Condition of cstorvolume. If underlying persistent volume is being
	// resized then the Condition will be set to 'ResizePending'.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []CStorVolumeCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	// ReplicaDetails refers to the trusty replica information which are
	// connected at given time
	ReplicaDetails CStorVolumeReplicaDetails `json:"replicaDetails,omitempty"`
}

// CStorVolumeReplicaDetails contains trusty replica inform which will be
// updated by target
type CStorVolumeReplicaDetails struct {
	// KnownReplicas represents the replicas that target can trust to read data
	KnownReplicas map[ReplicaID]string `json:"knownReplicas,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CStorVolumeList is a list of CStorVolume resources
type CStorVolumeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CStorVolume `json:"items"`
}

// CVStatusResponse stores the response of istgt replica command output
// It may contain several volumes
type CVStatusResponse struct {
	CVStatuses []CVStatus `json:"volumeStatus"`
}

// CVStatus stores the status of a CstorVolume obtained from response
type CVStatus struct {
	Name            string          `json:"name"`
	Status          string          `json:"status"`
	ReplicaStatuses []ReplicaStatus `json:"replicaStatus"`
}


// ReplicaStatus stores the status of replicas
type ReplicaStatus struct {
	// ID is replica unique identifier
	ID string `json:"replicaId"`
	// Mode represents replica status i.e. Healthy, Degraded
	Mode string `json:"mode"`
	// Represents IO number of replicas persisted on the disk
	CheckpointedIOSeq string `json:"checkpointedIOSeq"`
	// Ongoing reads I/O from target to replica
	InflightRead string `json:"inflightRead"`
	// ongoing writes I/O from target to replica
	InflightWrite string `json:"inflightWrite"`
	// Ongoing sync I/O from target to replica
	InflightSync string `json:"inflightSync"`
	// time since the replica connected to target
	UpTime int `json:"upTime"`
	// Quorum indicates wheather data wrtitten to the replica
	// is lost or exists.
	// "0" means: data has been lost( might be ephimeral case)
	// and will recostruct data from other Healthy replicas in a write-only
	// mode
	// 1 means: written data is exists on replica
	Quorum string `json:"quorum"`
}

// CStorVolumeCondition contains details about state of cstorvolume
type CStorVolumeCondition struct {
	// Type is a different valid value of CStorVolumeCondition
	Type   CStorVolumeConditionType `json:"type"`
	
	Status ConditionStatus `json:"status"`
	// Last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Unique, this should be a short, machine understandable string that gives the reason
	// for condition's last transition. If it reports "ResizePending" that means the underlying
	// cstorvolume is being resized.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// CStorVolumeConditionType is a valid value of CStorVolumeCondition.Type
type CStorVolumeConditionType string

const (
	// CStorVolumeResizing - a user trigger resize of pvc has been started
	CStorVolumeResizing CStorVolumeConditionType = "Resizing"
)

// ConditionStatus states in which state condition is present
type ConditionStatus string

// These are valid condition statuses. "ConditionInProgress" means corresponding
// condition is inprogress. "ConditionSuccess" means corresponding condition is success
const (
	// ConditionInProgress states resize of underlying volumes are in progress
	ConditionInProgress ConditionStatus = "InProgress"
	// ConditionSuccess states resizing underlying volumes are successful
	ConditionSuccess ConditionStatus = "Success"
)
```

## CStorVolumeReplicas APIs

```go

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=cstorvolumereplica

// CStorVolumeReplica describes a cstor volume resource created as custom resource
type CStorVolumeReplica struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CStorVolumeReplicaSpec   `json:"spec"`
	Status            CStorVolumeReplicaStatus `json:"status"`
	VersionDetails    VersionDetails           `json:"versionDetails"`
}


// CStorVolumeReplicaSpec is the spec for a CStorVolumeReplica resource
type CStorVolumeReplicaSpec struct {
	// TargetIP represents iscsi target IP through which replica cummunicates
	// IO workloads and other volume operations like snapshot and resize requests
	TargetIP string `json:"targetIP"`
	//Represents the actual capacity of the underlying volume
	Capacity string `json:"capacity"`
	// ZvolWorkers represents number of threads that executes client IOs
	ZvolWorkers string `json:"zvolWorkers"`
	// ReplicaID is unique number to identify the replica
	ReplicaID string `json:"replicaid"`
	// Controls the compression algorithm used for this volumes
	// examples: on|off|gzip|gzip-N|lz4|lzjb|zle
	Compression string `json:"compression"`
	// BlockSize is the logical block size in multiple of 512 bytes
	// BlockSize specifies the block size of the volume. The blocksize
	// cannot be changed once the volume has been written, so it should be
	// set at volume creation time. The default blocksize for volumes is 4 Kbytes.
	// Any power of 2 from 512 bytes to 128 Kbytes is valid.
	BlockSize uint32 `json:"blockSize"`
}


// CStorVolumeReplicaPhase is to hold result of action.
type CStorVolumeReplicaPhase string

// Status written onto CStorVolumeReplica objects.
const (

	// CVRStatusEmpty describes CVR resource is created but not yet monitored by
	// controller(i.e resource is just created)
	CVRStatusEmpty CStorVolumeReplicaPhase = ""

	// CVRStatusOnline describes volume replica is Healthy and data existing on
	// the healthy replica is up to date
	CVRStatusOnline CStorVolumeReplicaPhase = "Healthy"

	// CVRStatusOffline describes volume replica is created but not yet connected
	// to the target
	CVRStatusOffline CStorVolumeReplicaPhase = "Offline"

	// CVRStatusDegraded describes volume replica is connected to the target and
	// rebuilding from other replicas is not yet started but ready for serving
	// IO's
	CVRStatusDegraded CStorVolumeReplicaPhase = "Degraded"

	// CVRStatusNewReplicaDegraded describes replica is recreated (due to pool
	// recreation[underlying disk got changed]/volume replica scaleup cases) and
	// just connected to the target. Volume replica has to start reconstructing
	// entire data from another available healthy replica. Until volume replica
	// becomes healthy whatever data written to it is lost(NewReplica also not part
	// of any quorum decision)
	CVRStatusNewReplicaDegraded CStorVolumeReplicaPhase = "NewReplicaDegraded"

	// CVRStatusRebuilding describes volume replica has missing data and it
	// started rebuilding missing data from other replicas
	CVRStatusRebuilding CStorVolumeReplicaPhase = "Rebuilding"

	// CVRStatusReconstructingNewReplica describes volume replica is recreated
	// and it started reconstructing entire data from other healthy replica
	CVRStatusReconstructingNewReplica CStorVolumeReplicaPhase = "ReconstructingNewReplica"

	// CVRStatusError describes either volume replica is not exist in cstor pool
	CVRStatusError CStorVolumeReplicaPhase = "Error"

	// CVRStatusDeletionFailed describes volume replica deletion is failed
	CVRStatusDeletionFailed CStorVolumeReplicaPhase = "DeletionFailed"

	// CVRStatusInvalid ensures invalid resource(currently not honoring)
	CVRStatusInvalid CStorVolumeReplicaPhase = "Invalid"

	// CVRStatusInit describes CVR resource is newly created but it is not yet
	// created zfs dataset
	CVRStatusInit CStorVolumeReplicaPhase = "Init"

	// CVRStatusRecreate describes the volume replica is recreated due to pool
	// recreation/scaleup
	CVRStatusRecreate CStorVolumeReplicaPhase = "Recreate"
)

// CStorVolumeReplicaStatus is for handling status of cvr.
type CStorVolumeReplicaStatus struct {
	// CStorVolumeReplicaPhase is to holds different phases of replica
	Phase CStorVolumeReplicaPhase `json:"phase"`
	// CStorVolumeCapacityDetails represents capacity info of replica
	Capacity CStorVolumeReplicaCapacityDetails `json:"capacity"`
	// LastTransitionTime refers to the time when the phase changes
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The last updated time
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`

	// Snapshots contains list of snapshots, and their properties,
	// created on CVR
	Snapshots map[string]CStorSnapshotInfo `json:"snapshots,omitempty"`

	// PendingSnapshots contains list of pending snapshots that are not yet
	// available on this replica
	PendingSnapshots map[string]CStorSnapshotInfo `json:"pendingSnapshots,omitempty"`
}

// CStorSnapshotInfo represents the snapshot information related to particular
// snapshot
type CStorSnapshotInfo struct {
	// LogicalReferenced describes the amount of space that is "logically"
	// accessable by this snapshot. This logical space ignores the
	// effect of the compression and copies properties, giving a quantity
	// closer to the amount of data that application see. It also includes
	// space consumed by metadata.
	LogicalReferenced uint64 `json:"logicalReferenced"`

	// Used is the used bytes for given snapshot
	// Used uint64 `json:"used"`
}


// CStorVolumeReplicaCapacityDetails represents capacity information related to volume
// replica
type CStorVolumeReplicaCapacityDetails struct {
	// The amount of space consumed by this volume replica and all its descendants
	Total string `json:"total"`
	// The amount of space that is "logically" accessible by this dataset. The logical
	// space ignores the effect of the compression and copies properties, giving a
	// quantity closer to the amount of data that applications see.  However, it does
	// include space consumed by metadata
	Used string `json:"used"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=cstorvolumereplicas

// CStorVolumeReplicaList is a list of CStorVolumeReplica resources
type CStorVolumeReplicaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CStorVolumeReplica `json:"items"`
}

```
