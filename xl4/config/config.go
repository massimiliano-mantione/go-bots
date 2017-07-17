package config

const ForwardAcceleration = 10000 / 400
const ReverseAcceleration = 10000 / 1

const MaxIrValue = 100
const MaxIrDistance = 100

// VisionIntensityMax is the maximum vision intensity
const VisionIntensityMax = 100

// VisionAngleMax is the maximum vision angle (positive on the right)
const VisionAngleMax = 100

const MaxSpeed = 10000

const StartTime = 50

const OuterSpeed = MaxSpeed
const InnerSpeed = 4200
const AdjustInnerMax = 200

const SeekMoveSpeed = 4000
const SeekMoveMillis = 1000
const SeekTurnSpeed = 4000
const SeekTurnMillis = 1500

const BackTurn1SpeedOuter = MaxSpeed
const BackTurn1SpeedInner = MaxSpeed / 2
const BackTurn1Millis = 500
const BackTurn2Speed = 4000
const BackTurn2Millis = 1500
const BackMoveSpeed = MaxSpeed
const BackMoveMillis = 500