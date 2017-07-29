package config

const ColorIsOut = 5

const ForwardAcceleration = 10000 / 600
const ReverseAcceleration = 10000 / 1

//const MaxIrValue = 100

const MaxIrDistance = 80

const IgnoreBorderIrDistance = 40

// VisionIntensityMax is the maximum vision intensity
//const VisionIntensityMax = 100

// VisionAngleMax is the maximum vision angle (positive on the right)
//const VisionAngleMax = 100

const MaxSpeed = 10000

const StartTime = 500

const SeekMoveSpeed = 3500
const SeekMoveMillis = 860
const SeekTurnSpeed = 3700
const SeekTurnMillis = 1200

const BackTurn1SpeedOuter = MaxSpeed
const BackTurn1SpeedInner = MaxSpeed / 2
const BackTurn1Millis = 80
const BackTurn2Speed = 3700
const BackTurn2Millis = 500
const BackMoveSpeed = MaxSpeed
const BackMoveMillis = 5
const BackTurn3Speed = 3700
const BackTurn3Millis = 1600

const CircleFindBorderMillis = 150
const CircleFindBorderOuterSpeed = MaxSpeed * 80 / 100
const CircleFindBorderInnerSpeed = MaxSpeed * 40 / 100
const CircleFindBorderOuterSpeedSlowLeft = MaxSpeed * 28 / 100
const CircleFindBorderInnerSpeedSlowLeft = MaxSpeed * 18 / 100
const CircleFindBorderOuterSpeedSlowRight = MaxSpeed * 28 / 100
const CircleFindBorderInnerSpeedSlowRight = MaxSpeed * 18 / 100
const CircleMillis = 2500
const CircleOuterSpeed = MaxSpeed
const CircleInnerSpeedLeft = 3000
const CircleInnerSpeedRight = 3000
const CircleAdjustInnerMax = 500
const CircleSpiralMillis = 450
const CircleSpiralOuterSpeed = MaxSpeed
const CircleSpiralInnerSpeed = 1500

const GoForwardMillis = 700
const GoForwardSpeed = MaxSpeed

// const GoForwardTurnMillis = 600
const GoForwardTurnMillis = 0
const GoForwardTurnOuterSpeed = MaxSpeed
const GoForwardTurnInnerSpeed = 1000

const TurnBackPreMoveMillis = 400
const TurnBackPreMoveSpeed = MaxSpeed
const TurnBackMillis = 120
const TurnBackOuterSpeed = MaxSpeed
const TurnBackInnerSpeed = -MaxSpeed
const TurnBackMoveMillis = 800
const TurnBackMoveSpeed = MaxSpeed

const TrackOnly1SensorOuterSpeed = MaxSpeed
const TrackOnly1SensorInnerSpeed = 8000
const TrackSpeed = MaxSpeed
const TrackCenterZone = 10
const TrackDifferenceCoefficent = 50
