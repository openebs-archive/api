---
oep-number: 2
title: Pool Expansion When Underlying Disk Expanded
authors:
  - "@mittachaitu"
owners:
  - "@vishnuitta"
  - "@kmova"
  - "@prateekpandey14"
  - "@sonasingh46"
editor: "@mittachaitu"
creation-date: 2020-03-20
last-updated: 2020-03-20
status: Implementable
---

# Pool Expansion When Underlying Disk Expanded

## Table of Contents

- [Pool Expansion When Underlying Disk Expanded](#pool-expansion-when-underlying-disk-expanded)
	- [Table of Contents](#table-of-contents)
	- [Summary](#summary)
	- [Motivation](#motivation)
		- [Goals](#goals)
		- [Non-Goals](#non-goals)
	- [Proposal](#proposal)
		- [User Stories](#user-stories)
		- [Current Implementation](#current-implementation)
		- [Proposed Implementations](#proposed-implementations)
		- [Steps to perform user stories](#steps-to-perform-user-stories)
			- [Pool Expansion Steps:](#pool-expansion-steps)
		- [Low Level Design](#low-level-design)
			- [Work Flow](#work-flow)
	- [Drawbacks](#drawbacks)
	- [Alternatives](#alternatives)
	- [Infrastructure Needed](#infrastructure-needed)

## Summary

- This proposal includes design details of the pool expansion when underlying disks were expanded.

## Motivation

There were cases where OpenEBS users expand the pool by expanding underlying disks(which were used by pool) when space on the pool got exhausted.

### Goals

- Expand cStor pool when the underlying disk got expanded.

### Non-Goals

- If user dedicated the first part of partitioned disk for cStor and second part for third-party then pool expansion can't be performed.

## Proposal

### User Stories

As a OpenEBS user I should be able to expand the pool when underlying disk got expanded.

### Current Implementation

- When user creates the CSPC, CSPC controller which is watching for CSPC will get create request and creates cstorpoolinstances and cstorpooldeployments. CstorPoolInstance is responsible for holding the blockdevice information and other pool configurations and cStor pool deployment have 3 containers cstor-pool-mgmt, cstor-pool and m-exporter as name conveyes cstor-pool container creates pool(ZFS file system) on top of the blockdevices, where as cstor-pool-mgmt sidecar container is responsible for managing these pool. If any changes required as of day 2 operations on the pool then pool configurations will be updated by adding BlockDevices, adding RaidGroup and replacing blockdevices cstor-pool-mgmt looks for any pending pool operations for CStorPoolInstance then performs requested operation on the pool.

Below are high level info to explain the Day2 operations that are currently supported.

- Pool Expansion When BlockDevice/raid group added:
	cstor-pool-mgmt will iterate over all the BlockDevices that exists in CStorPoolInstance resource and verify whether the corresponding BlockDevice/raid group is part of the pool. If that BlockDevice/raid group is not part of the pool then pool expansion will be triggered on newly added BlockDevice/raid group.

- BlockDevice Replacement In Pool:
    cstor-pool manager will iterate over all the BlockDevices present in CstorPoolInstance and verifies the requested BlockDevice replacement request (we are maintaining annotation called openebs.io/predecessor: <replaced_blockdevice_name> on the claim of that blockdevice). If requested BlockDevice is valid for replacement then the pool managet will perform BlockDevice replacement operation on the pool.

### Proposed Implementations

- CStor-pool Manger inbuilt CSPI controller watches for CSPI resources and tries to get the desired state via every sync operation. It will iterate over all the BlockDevices that exist on CStorPoolInstance resource and if BlockDevice resource capacity(considerable changes) is greater than the existing capacity of BlockDevice on CstorPoolInstanceBlockDevice then we can follow below-mentioned steps to perform the pool expansion:
1. If IsBlockDeviceExpanded flag sets true then it is something like user/admin is requesting to perform pool expansion steps.
2. If IsBlockDeviceExpanded flag sets false then NoOps will be performed.
[What if cStor is consuming part of partitioned disk? User has to set IsBlockDeviceExpanded to true and other part of partitioned disk shouldn't have any FileSystem then only pool expansion will be triggered].
[What if cStor is consuming entier partitioned disk? User has to set IsBlockDeviceExpanded to true then cstor-pool-mgmt will expand the partition and trigger pool expansion process].

Currently CstorPoolInstanceBlockDevice holds below configurations 
```go
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
```

To get user Input proposed to add field in above struct

```go
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

	// IsBlockDeviceExpanded informs whether blockdevice 
	// got expanded or not
	IsBlockDeviceExpanded bool `json:”isBlockDeviceExpanded”`
}
```

### Steps to perform user stories
- User has to set the value of IsBlockDeviceExpanded to true in corresponding CStorPoolInstanceBlockDevice.
- Once the operation is completed user can unset the value of ISBlockDeviceExpanded[How users can identify
  operation is succeeded? User can know it via events (or) Increased capacity in CStorPoolInstanceBlockDevice].
  NOTE: Leaving IsBlockDeviceExpanded can be left set it will not cause any extra network calls.

#### Pool Expansion Steps:

Step1: Set autoexpand to on in pool by triggering `zpool set autoexpand=on <pool_name>`.

Step2: Remove the buffering partition if exists which was created by cstor-pool container(ex: sdb9).

Step3: Expand the existing partition by executing following comamnd
       ```sh
       parted </disk> resizepart <partition_number> 100%
       ```

Step4: Inform the pool saying disk was expanded by triggering following command
       ```sh
       zpool online -e <pool_name> <disk_name>
       ```
NOTE: Above steps automates the steps mentioned in [doc](https://github.com/openebs/openebs-docs/blob/day_2_ops/docs/resize-single-disk-pool.md).

### Low Level Design

#### Work Flow

- User/Admin will set the value of IsBlockDeviceExpanded on CSPC to true to expand the cstor pool which was on top of blockdevice.
- Once User/Admin updates the value CSPC controller will get update event. As part of event processing CSPC controller will identify which  CSPI needs to be updated and it will update CSPI with updated information.
- When CSPI is updated CSPI controller will get update event and verifies which disk was expanded and performs Pool Expansion Steps only if IsBlockDeviceExpanded holds true and capacity of on blockdevice CR is greater than capacity on CStorPoolInstanceBlockDevice[What if only one blockdevice got expanded in case of raid/mirror configuration? Pool expansion operation will be triggered but there won't be any expansion on pool].

NOTE: Pool will be expanded only if all the blockdevices in raidgroup were expanded in case of mirror/raid configuration. If it is stripe configuration expanding one blockdevice will expand the pool.

## Drawbacks

NA

## Alternatives

NA

## Infrastructure Needed

NA
