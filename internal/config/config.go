package config

var Protocol = "tcp"
var Port = ":3000"
var MaxConnection = 20000

var EvictionMaxKeyNumber = 10
var EvictionRation = 0.1
var EvictionPolicy = "allkeys-random"

var EpoolMaxSize = 16
var EpoolLruSampleSize = 5

var AvgTtlRandomSampleSize = 5
