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

NA

## Proposal

### User Stories

As a OpenEBS user I should be able to expand the pool when underlying disk got expanded.

### Current Implementation

When user creates the CSPC, CSPC controller which is watching for CSPC will get create request and creats cstorpoolinstances and cstorpooldeployments. CstorPoolInstance is responsible for holding the blockdevice information and other pool configurations while cStor pool deployment contains following containers cstor-pool-mgmt, cstor-pool container and m-exporter as name conveyes cstor-pool-mgmt is responsible for managing the pool. cstor-pool container creates pool(ZFS file system) on top of the blockdevices. If any changes need to make on pool configuration like adding blockdevices, adding raid groups and replacing blockdevice cstor-pool-mgmt container will look into CstorPoolInstance for pool information and verifies are there any pending actions to be performed then cstor-pool-mgmt will prforms corresponding operation on the pool. Below info is hight level info to know how the operation is detected

- Pool Expansion When blockdevice/raid group added:	
	cstor-pool-mgmt will iterate over all the blockdevices present in CstorPoolInstance and then verifies whether corresponding blockdevice is in use in the pool. If that blockdevice/raid group doesn’t exist then trigger is for pool expansion by adding blockdevice/raid group.

- BlockDevice replacement:
	cstor-pool-mgmt will iterate over all the blockdevices present in CstorPoolInstance and the verifies is that blockdevice is for replacement(We are maintaining annotation called openebs.io/predecessor: <replaced_blockdevice_name> on the claim of that blockdevice) if that block device is for replacement then the it will perform replacement operation on the pool.

### Proposed Implementations

CStor-pool-mgmt has CSPI controller for every sync operation cstor-pool-mgmt will iterate over all the blockdevices exist on CSPI by iterating over CstorPoolInstanceBlockDevices of all raid groups. Now if the capacity exist on blockdevice CR is greater than the  capacity exists on CstorPoolInstanceBlockDevice then we can follow the steps to perform the pool expansion. But anyways it is not good practice to perform expansion of pool with out getting the permission from user. So to get the input from the user we will be introducing a new field under CstorPoolInstanceBlockDevice called IsBlockDeviceExpanded. If  IsBlockDeviceExpanded sets true then it is something like user/admin is requesting to perform pool expansion steps. If IsBlockDeviceExpanded sets false then NoOps will be performed.

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

#### Pool Expansion Steps:

Step1: Set autoexpand to on in pool by triggering `zpool set autoexpand=on <pool_name>`.

Step2: Remove the partition if any exists which was created by cstor-pool container(ex: sdb9).

Step3: Expand the existing partition by executing following comamnd (parted </disk> resizepart <partition_number> 100%).

Step4: Inform the pool saying disk was expanded by triggering following command zpool online -e <pool_name> <disk_name>.


### Low Level Design

#### Work Flow
- User/Admin will set the value of IsBlockDeviceExpanded to true to expand the cstor pool which was on top of blockdevice.
- Once User/Admin updates the value CSPC controller will get update event. As part of event processing CSPC controller will identify which  CSPI needs to be updated and it will update CSPI with updated information.
- When CSPI is updated CSPI controller will get update event and verifies which disk was expanded and performs Pool Expansion Steps only if IsBlockDeviceExpanded holds true and capacity of on blockdevice CR is greater than capacity on CStorPoolInstanceBlockDevice.

## Drawbacks

NA

## Alternatives

NA

## Infrastructure Needed

NA
