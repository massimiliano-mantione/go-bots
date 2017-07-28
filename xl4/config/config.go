package config

const ColorIsOut = 5

const ForwardAcceleration = 10000 / 400
const ReverseAcceleration = 10000 / 1

//const MaxIrValue = 100

const MaxIrDistance = 80
const IgnoreBorderIrDistance = 40

// VisionIntensityMax is the maximum vision intensity
//const VisionIntensityMax = 100

// VisionAngleMax is the maximum vision angle (positive on the right)
//const VisionAngleMax = 100

const MaxSpeed = 10000

const StartTime = 5000

const SeekMoveSpeed = 3500
const SeekMoveMillis = 860
const SeekTurnSpeed = 4000
const SeekTurnMillis = 1100

const BackTurn1SpeedOuter = MaxSpeed
const BackTurn1SpeedInner = MaxSpeed / 2
const BackTurn1Millis = 400
const BackTurn2Speed = 4000
const BackTurn2Millis = 300
const BackMoveSpeed = MaxSpeed
const BackMoveMillis = 200
const BackTurn3Speed = 4000
const BackTurn3Millis = 800

const CircleFindBorderMillis = 250
const CircleFindBorderOuterSpeed = MaxSpeed * 80 / 100
const CircleFindBorderInnerSpeed = MaxSpeed * 40 / 100
const CircleFindBorderOuterSpeedSlowLeft = MaxSpeed * 25 / 100
const CircleFindBorderInnerSpeedSlowLeft = MaxSpeed * 18 / 100
const CircleFindBorderOuterSpeedSlowRight = MaxSpeed * 30 / 100
const CircleFindBorderInnerSpeedSlowRight = MaxSpeed * 20 / 100
const CircleMillis = 1700
const CircleOuterSpeed = MaxSpeed
const CircleInnerSpeedLeft = 3500
const CircleInnerSpeedRight = 3500
const CircleAdjustInnerMax = 400
const CircleSpiralMillis = 450
const CircleSpiralOuterSpeed = MaxSpeed
const CircleSpiralInnerSpeed = 2000

const GoForwardMillis = 600
const GoForwardSpeed = MaxSpeed

// const GoForwardTurnMillis = 600
const GoForwardTurnMillis = 0
const GoForwardTurnOuterSpeed = MaxSpeed
const GoForwardTurnInnerSpeed = 1000

const TurnBackPreMoveMillis = 300
const TurnBackPreMoveSpeed = MaxSpeed
const TurnBackMillis = 180
const TurnBackOuterSpeed = MaxSpeed
const TurnBackInnerSpeed = -MaxSpeed
const TurnBackMoveMillis = 300
const TurnBackMoveSpeed = MaxSpeed

const TrackOnly1SensorOuterSpeed = MaxSpeed
const TrackOnly1SensorInnerSpeed = 8000
const TrackSpeed = MaxSpeed
const TrackCenterZone = 10
const TrackDifferenceCoefficent = 50
