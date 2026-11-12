package security

import "time"

// timeNow is a small seam so tests can pin TOTP time if they need to.
var timeNow = time.Now
