---
page_title: "genesyscloud_outbound_settings Resource - terraform-provider-genesyscloud"
subcategory: ""
description: |-
  An organization's outbound settings
---
# genesyscloud_outbound_settings (Resource)

An organization's outbound settings

## API Usage
The following Genesys Cloud APIs are used by this resource. Ensure your OAuth Client has been granted the necessary scopes and permissions to perform these operations:

* [GET /api/v2/outbound/settings](https://developer.genesys.cloud/routing/outbound/#get-api-v2-outbound-settings)
* [PATCH /api/v2/outbound/settings](https://developer.genesys.cloud/routing/outbound/#patch-api-v2-outbound-settings)

## Example Usage

```terraform
resource "genesyscloud_outbound_settings" "example_settings" {
  max_calls_per_agent                 = 10
  max_line_utilization                = 0.5
  abandon_seconds                     = 6.5
  compliance_abandon_rate_denominator = "ALL_CALLS"
  automatic_time_zone_mapping {
    callable_windows {
      mapped {
        earliest_callable_time = "09:00"
        latest_callable_time   = "17:00"
      }
      unmapped {
        earliest_callable_time = "08:00"
        latest_callable_time   = "18:00"
        time_zone_id           = "CET"
      }
    }
    supported_countries = ["US"]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `abandon_seconds` (Number) The number of seconds used to determine if a call is abandoned.
- `automatic_time_zone_mapping` (Block List) The settings for automatic time zone mapping. Note that changing these settings will change them for both voice and messaging campaigns. (see [below for nested schema](#nestedblock--automatic_time_zone_mapping))
- `compliance_abandon_rate_denominator` (String) The denominator to be used in determining the compliance abandon rate.Valid values: ALL_CALLS, CALLS_THAT_REACHED_QUEUE.
- `max_calls_per_agent` (Number) The maximum number of calls that can be placed per agent on any campaign.
- `max_line_utilization` (Number) The maximum percentage of lines that should be used for Outbound, expressed as a decimal in the range [0.0, 1.0].
- `reschedule_time_zone_skipped_contacts` (Boolean) Whether or not to reschedule time-zone blocked contacts.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--automatic_time_zone_mapping"></a>
### Nested Schema for `automatic_time_zone_mapping`

Optional:

- `callable_windows` (Block Set, Max: 1) The time intervals to use for automatic time zone mapping. (see [below for nested schema](#nestedblock--automatic_time_zone_mapping--callable_windows))
- `supported_countries` (List of String) The countries that are supported for automatic time zone mapping.

<a id="nestedblock--automatic_time_zone_mapping--callable_windows"></a>
### Nested Schema for `automatic_time_zone_mapping.callable_windows`

Optional:

- `mapped` (Block Set, Max: 1) The time interval to place outbound calls, for contacts that can be mapped to a time zone. (see [below for nested schema](#nestedblock--automatic_time_zone_mapping--callable_windows--mapped))
- `unmapped` (Block Set, Max: 1) The time interval and time zone to place outbound calls, for contacts that cannot be mapped to a time zone. (see [below for nested schema](#nestedblock--automatic_time_zone_mapping--callable_windows--unmapped))

<a id="nestedblock--automatic_time_zone_mapping--callable_windows--mapped"></a>
### Nested Schema for `automatic_time_zone_mapping.callable_windows.mapped`

Optional:

- `earliest_callable_time` (String) The earliest time to dial a contact. Valid format is HH:mm
- `latest_callable_time` (String) The latest time to dial a contact. Valid format is HH:mm.


<a id="nestedblock--automatic_time_zone_mapping--callable_windows--unmapped"></a>
### Nested Schema for `automatic_time_zone_mapping.callable_windows.unmapped`

Optional:

- `earliest_callable_time` (String) The earliest time to dial a contact. Valid format is HH:mm.
- `latest_callable_time` (String) The latest time to dial a contact. Valid format is HH:mm.
- `time_zone_id` (String) The time zone to use for contacts that cannot be mapped.

