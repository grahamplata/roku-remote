package roku

// UnableToLocate User instructions if the cli is unable to locate any devices
const UnableToLocate string = `Unable to locate any Roku Devices on your network.
Please check if your device is online and connected to your network.`

// NoDefaultRoku User instructions if there is no default roku assigned
const NoDefaultRoku string = "Consider running roku-remote find command first to set a default device"

// MissingAction User instructions if there is not a provided action
const MissingAction string = `You need to include an action Try using --help`
