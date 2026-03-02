# Feature Flags

## Overview

The Feature Flags module provides creation and management of feature flags within a space. Flags can be boolean or multivariate (multiple variations with distinct values) and can be toggled on or off at any time.

## Core Concepts

### `FeatureFlag`

| Field | Type | Description |
|-------|------|-------------|
| `ID` | int64 | Primary key |
| `SpaceID` | int64 | FK to spaces |
| `Identifier` | string | Space-unique identifier |
| `Name` | string | Display name |
| `Description` | string | Optional description |
| `Kind` | string | Flag kind (e.g., `boolean`, `multivariate`) |
| `DefaultOnVariation` | string | Identifier of the variation served when flag is on |
| `DefaultOffVariation` | string | Identifier of the variation served when flag is off |
| `Enabled` | bool | Whether the flag is currently on |
| `Variations` | []Variation | List of variations |
| `Tags` | []string | Tag array |
| `Permanent` | bool | Whether the flag is permanent (not subject to cleanup) |
| `CreatedBy` | int64 | FK to principals |
| `Created` | int64 | Unix milliseconds |
| `Updated` | int64 | Unix milliseconds |
| `Version` | int64 | Optimistic locking version |

### `Variation`

| Field | Description |
|-------|-------------|
| `Identifier` | Unique key for this variation |
| `Name` | Display name |
| `Value` | The variation value (any type) |

## API Endpoints

### Create Flag

**POST** `/api/v1/spaces/{space_ref}/feature-flags`

```json
{
  "identifier": "new-checkout-flow",
  "name": "New Checkout Flow",
  "description": "Enables the redesigned checkout experience",
  "kind": "boolean",
  "default_on_variation": "true",
  "default_off_variation": "false",
  "enabled": false,
  "variations": [
    {"identifier": "true",  "name": "Enabled",  "value": true},
    {"identifier": "false", "name": "Disabled", "value": false}
  ],
  "tags": ["checkout", "ui"],
  "permanent": false
}
```

### List Flags

**GET** `/api/v1/spaces/{space_ref}/feature-flags`

Query parameters: `query` (search), standard pagination

### Get Flag

**GET** `/api/v1/spaces/{space_ref}/feature-flags/{identifier}`

### Update Flag

**PATCH** `/api/v1/spaces/{space_ref}/feature-flags/{identifier}`

Fields that can be updated include `name`, `description`, `enabled`, `variations`, `tags`, `default_on_variation`, `default_off_variation`.

### Delete Flag

**DELETE** `/api/v1/spaces/{space_ref}/feature-flags/{identifier}`

## MCP Integration

The `feature_flag_toggle` atomic MCP tool wraps the feature flag controller, allowing AI agents to toggle flags via the MCP protocol.

## Web UI

The **Feature Flags** page (`/pages/FeatureFlagList`) and the **Feature Flags** navigation item in the module sidebar provide a web interface for managing flags.

## File Locations

| Purpose | Path |
|---------|------|
| Types | `types/featureflag.go` |
| Controller | `app/api/controller/featureflag/controller.go` |
| Create | `app/api/controller/featureflag/create.go` |
| Update | `app/api/controller/featureflag/update.go` |
| Handlers | `app/api/handler/featureflag/` |
