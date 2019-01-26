# ovaify

## What is it?
A Go library for creating OVA (Open Virtual Appliance) files.

## Use cases
This library was developed to simplify the virtual machine deployment supply
chain. It enables the creation of OVA (Open Virtual Appliance) files from
an existing OVF (Open Virtualization Format) file and its associated artifacts
(such as a virtual machine disk image).

An OVF only provides the configuration for a virtual machine appliance - it
does not provide the appliance's disk, or other files. A OVA on the other hand
provides all of these in the form of a single compressed file. Using OVAs makes
deploying new appliances easier, and more maintainable.

While open source tools like [packer](https://packer.io) and
[VirtualBox](https://www.virtualbox.org/) can create these files, they cannot
easily create OVA files from existing OVFs. This is usually worked around using
VMWare's [ovftool](https://www.vmware.com/support/developer/ovf/). Because
ovftool is closed source, incorporating it into a VM development toolchain can
be a logistical headache. This library allows developers to incorporate
ovftool's functionality into their toolchain without such headaches.

## API
The library's API is very small. The most notable function is the
`CreateOvaFile` function. This function creats an OVA using the provided
`OvaConfig`. Here is an example application that uses this function:
```go
package main

import (
    "log"

    "github.com/stephen-fox/ovaify"
)

func main() {
    config := ovaify.OvaConfig{
        OutputFilePath:     "/my-awesome.ova",
        OvfFilePath:        "/my-vm.ovf",
        FilePathsToInclude: []string{
            "/my-vm-disk-image.vmdk",
        },
    }

    err := ovaify.CreateOvaFile(config)
    if err != nil {
        log.Fatal("Failed to create OVA - " + err.Error())
    }
}
```
