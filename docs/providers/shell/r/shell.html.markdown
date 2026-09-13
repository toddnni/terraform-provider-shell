---
layout: "script"
page_title: "Script: script_wrapper"
sidebar_current: "docs-script-resource-wrapper"
description: |-
  Provides a PowerDNS record resource.
---

# powerdns\_record

Provides a PowerDNS record resource.

## Example Usage

```
# Add a record to the zone
resource "powerdns_record" "foobar" {
	zone = "example.com"
	name = "www.example.com"
	type = "A"
	ttl = 300
	records = ["192.168.0.11"]
}
```

## Argument Reference

The following arguments are supported:

* `zone` - (Required) The name of zone to contain this record.
* `name` - (Required) The name of the record.
* `type` - (Required) The record type.
* `ttl` - (Required) The TTL of the record.
* `records` - (Required) A string list of records.

