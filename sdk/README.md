# One-Log SDKs

Official SDKs for One-Log (ULAM) - Unified Log & Activity Monitor

## Available SDKs

- **Go** - `sdk/go/` ✅ Complete
- **Node.js** - `sdk/node/` ✅ Complete  
- **Python** - `sdk/python/` ✅ Complete
- **PHP** - `sdk/php/` ✅ Complete

## Quick Start

### Go

```go
import "github.com/petrushandika/one-log/sdk/go/onelog"

client := onelog.NewClient(
    onelog.WithAPIKey("ulam_live_xxx"),
)

client.LogInfo(onelog.SystemError, "Server started", nil)
```

### Node.js

```javascript
const OneLog = require('@onelog/sdk');

const client = new OneLog({
  apiKey: 'ulam_live_xxx'
});

await client.logInfo('SYSTEM_ERROR', 'Server started');
```

### Python

```python
from onelog import OneLog

client = OneLog(api_key="ulam_live_xxx")
client.log_info("SYSTEM_ERROR", "Server started")
```

### PHP

```php
use OneLog\OneLog;

$client = new OneLog(['api_key' => 'ulam_live_xxx']);
$client->logInfo(OneLog::CATEGORY_SYSTEM_ERROR, 'Server started');
```

## Features

All SDKs support:

- ✅ Sync logging (blocking)
- ✅ Async logging (fire-and-forget)
- ✅ Performance metrics logging
- ✅ Authentication events
- ✅ Audit trail events
- ✅ Error logging with stack traces
- ✅ Middleware support (Express, Flask, etc.)
- ✅ Environment variable configuration
- ✅ Custom timeout settings

## Installation

See individual SDK README files for installation instructions.

## License

MIT License
