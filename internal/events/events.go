package events

// SystemEvent represents a system-level event.
type SystemEvent int

const (
    SystemEventNone SystemEvent = iota
    LaptopSuspend
    LaptopResume
    MicMuteLedOn
    MicMuteLedOff
    MicMuteLedToggle
    BacklightOff
    BacklightLow
    BacklightMedium
    BacklightHigh
    BacklightToggle
    SecondaryDisplayToggle
    USBKeyboardAttached
    USBKeyboardDetached
    TouchpadDisable
    TouchpadEnable
    TouchpadToggle
)

// String returns a human‑readable name for the system event.
func (s SystemEvent) String() string {
    switch s {
    case SystemEventNone:
        return "None"
    case LaptopSuspend:
        return "LaptopSuspend"
    case LaptopResume:
        return "LaptopResume"
    case MicMuteLedOn:
        return "MicMuteLedOn"
    case MicMuteLedOff:
        return "MicMuteLedOff"
    case MicMuteLedToggle:
        return "MicMuteLedToggle"
    case BacklightOff:
        return "BacklightOff"
    case BacklightLow:
        return "BacklightLow"
    case BacklightMedium:
        return "BacklightMedium"
    case BacklightHigh:
        return "BacklightHigh"
    case BacklightToggle:
        return "BacklightToggle"
    case SecondaryDisplayToggle:
        return "SecondaryDisplayToggle"
    case USBKeyboardAttached:
        return "USBKeyboardAttached"
    case USBKeyboardDetached:
        return "USBKeyboardDetached"
    case TouchpadDisable:
        return "TouchpadDisable"
    case TouchpadEnable:
        return "TouchpadEnable"
    case TouchpadToggle:
        return "TouchpadToggle"
    default:
        return "Unknown"
    }
}
