package np2p_const

import "time"

const ServerImplVersion uint16 = 1
const PacketStructureVersion uint16 = 3

const PubkeySize = 32
const EventIdSize = 32
const SignatureSize = 64
const ResendCcheckInterval = time.Minute * 1
const ResendTimeBaseMin = 5
const ResendMaxTimes = 10                            // Max time is 5*2^10 = 5120 minutes = about 3.5 days
const MemoryUsageLimitForDBBuffer = 50 * 1024 * 1024 // 50MB
const DBStoreEventDataNumMax = 100 * 1024            // about 50MB (assume 500bytes/event)

const ProfileAndFollowDataUpdateCheckIntervalSec = 60 * 60 * 24 // 1day
const NoResendReqSendIntervalSec = 60 * 60                      // 1hour
