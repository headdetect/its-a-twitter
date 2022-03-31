package model

// TODO: Talk about how this is just a hack. Would normally use something like Redis for this

// Temporary memory stores //
var Sessions map[string]*User = make(map[string]*User) // [authToken] = user
var Timelines map[string]*GeneratedTimeline = make(map[string]*GeneratedTimeline) // [username] = timeline
