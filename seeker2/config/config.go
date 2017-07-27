package config

const ColorIsOut = 10

const ForwardAcceleration = 10000 / 200
const ReverseAcceleration = 10000 / 1

const MaxSpeed = 10000

const FrontWheelsSpeed = 100

const StartTime = 2000

const SeekMoveSpeed = 6000
const SeekMoveMillis = 850
const SeekTurnSpeed = 5000
const SeekTurnMillis = 1200

const BackTurn1SpeedOuter = MaxSpeed
const BackTurn1SpeedInner = MaxSpeed / 2
const BackTurn1Millis = 400
const BackTurn2Speed = 5000
const BackTurn2Millis = 800
const BackMoveSpeed = MaxSpeed
const BackMoveMillis = 500
const BackTurn3Speed = 5000
const BackTurn3Millis = 1000

const CircleFindBorderMillis = 300
const CircleFindBorderOuterSpeed = MaxSpeed * 60 / 100
const CircleFindBorderInnerSpeed = MaxSpeed * 20 / 100
const CircleFindBorderOuterSpeedSlow = MaxSpeed * 50 / 100
const CircleFindBorderInnerSpeedSlow = MaxSpeed * 30 / 100
const CircleMillis = 2000
const CircleOuterSpeed = MaxSpeed

// const CircleInnerSpeedLeft = 3700
const CircleInnerSpeedLeft = 4100

// const CircleInnerSpeedRight = 3800
const CircleInnerSpeedRight = 4200
const CircleAdjustInnerMax = 400
const CircleSpiralMillis = 500
const CircleSpiralOuterSpeed = MaxSpeed
const CircleSpiralInnerSpeed = 2000

const GoForwardMillis = 1000
const GoForwardSpeedLeft = MaxSpeed
const GoForwardSpeedRight = MaxSpeed

// const GoForwardTurnMillis = 600
const GoForwardTurnMillis = 0
const GoForwardTurnOuterSpeed = MaxSpeed
const GoForwardTurnInnerSpeed = 1000

const TurnBackMillis = 200
const TurnBackOuterSpeed = MaxSpeed
const TurnBackInnerSpeed = -MaxSpeed
const TurnBackMoveMillis = 500
const TurnBackMoveSpeed = MaxSpeed

const TrackFrontAngle = 30
const TrackMaxSpeed = MaxSpeed
const TrackInnerSpeed = 0
const TrackSpeedReductionAngle = VisionMaxAngle - TrackFrontAngle
const TrackSpeedReductionMax = TrackMaxSpeed - TrackInnerSpeed
const TrackVisionIntensityIgnoreBorder = 40

// VisionSpeed is the speed of the eyes motor, max is 1560
const VisionSpeed = "500"

// const VisionFarValueFront = 90
// const VisionFarValueSide = 80
const VisionFarValueFront = 70
const VisionFarValueSide = 60
const VisionFarValueDelta = VisionFarValueFront - VisionFarValueSide

const VisionMaxValue = 100
const VisionMaxAngle = 150
const VisionStartPosition = 150
const VisionStartPositionString = "150"
const VisionMaxPosition = 158
const VisionThresholdPosition = 155
const VisionEstimateReductionRange = 10
const VisionSpotWidth = 20
const VisionSpotSearchWidth = VisionMaxPosition - VisionSpotWidth
