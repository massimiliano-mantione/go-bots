package config

const ColorIsOut = 10

const ForwardAcceleration = 10000 / 200
const ReverseAcceleration = 10000 / 1

//const MaxIrValue = 100

const MaxIrDistance = 100

// VisionIntensityMax is the maximum vision intensity
//const VisionIntensityMax = 100

// VisionAngleMax is the maximum vision angle (positive on the right)
//const VisionAngleMax = 100

const MaxSpeed = 10000

const FrontWheelsSpeed = MaxSpeed

const StartTime = 20

const SeekMoveSpeed = 6000
const SeekMoveMillis = 900
const SeekTurnSpeed = 7500
const SeekTurnMillis = 900

const BackTurn1SpeedOuter = MaxSpeed
const BackTurn1SpeedInner = MaxSpeed / 2
const BackTurn1Millis = 400
const BackTurn2Speed = 7000
const BackTurn2Millis = 400
const BackMoveSpeed = MaxSpeed
const BackMoveMillis = 400
const BackTurn3Speed = 7000
const BackTurn3Millis = 800

const CircleFindBorderMillis = 500
const CircleFindBorderOuterSpeed = MaxSpeed * 90 / 100
const CircleFindBorderInnerSpeed = MaxSpeed * 50 / 100
const CircleFindBorderOuterSpeedSlow = MaxSpeed * 50 / 100
const CircleFindBorderInnerSpeedSlow = MaxSpeed * 30 / 100
const CircleMillis = 1700
const CircleOuterSpeed = MaxSpeed
const CircleInnerSpeedLeft = 3700
const CircleInnerSpeedRight = 3800
const CircleAdjustInnerMax = 300
const CircleSpiralMillis = 450
const CircleSpiralOuterSpeed = MaxSpeed
const CircleSpiralInnerSpeed = 2000

const GoForwardMillis = 350
const GoForwardSpeedLeft = 10000
const GoForwardSpeedRight = 7000

// const GoForwardTurnMillis = 600
const GoForwardTurnMillis = 0
const GoForwardTurnOuterSpeed = MaxSpeed
const GoForwardTurnInnerSpeed = 1000

const TurnBackMillis = 60
const TurnBackOuterSpeed = MaxSpeed
const TurnBackInnerSpeed = -MaxSpeed
const TurnBackMoveMillis = 300
const TurnBackMoveSpeed = MaxSpeed

const TrackOnly1SensorOuterSpeed = MaxSpeed
const TrackOnly1SensorInnerSpeed = 8000
const TrackSpeed = MaxSpeed
const TrackCenterZone = 10
const TrackDifferenceCoefficent = 50

// VisionSpeed is the speed of the eyes motor, max is 1560
const VisionSpeed = "900"
const VisionMaxPosition = 158
const VisionThresholdPosition = 155
