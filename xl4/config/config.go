package config

const ColorIsOut = 10

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
const BackMoveMillis = 500
const BackTurn3Speed = 4000
const BackTurn3Millis = 800

const CircleFindBorderMillis = 90
const CircleFindBorderOuterSpeed = MaxSpeed * 90 / 100
const CircleFindBorderInnerSpeed = MaxSpeed * 50 / 100
const CircleFindBorderOuterSpeedSlow = MaxSpeed * 30 / 100
const CircleFindBorderInnerSpeedSlow = MaxSpeed * 20 / 100
const CircleMillis = 150
const CircleOuterSpeed = MaxSpeed
const CircleInnerSpeed = 3600
const CircleAdjustInnerMax = 300
const CircleSpiralMillis = 20
const CircleSpiralOuterSpeed = MaxSpeed
const CircleSpiralInnerSpeed = 2500

const GoForwardMillis = 600
const GoForwardSpeed = MaxSpeed
const GoForwardTurnMillis = 20
const GoForwardTurnOuterSpeed = MaxSpeed
const GoForwardTurnInnerSpeed = 1000

const TurnBackMillis = 70
const TurnBackOuterSpeed = MaxSpeed
const TurnBackInnerSpeed = -MaxSpeed
const TurnBackMoveMillis = 500
const TurnBackMoveSpeed = MaxSpeed
