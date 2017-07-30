package config

const ColorIsOut = 20

const ForwardAcceleration = 10000 / 200
const ReverseAcceleration = 10000 / 1

const MaxSpeed = 10000

const FrontWheelsSpeed = 100

const StartTime = 5000

const SeekMoveSpeed = 6000
const SeekMoveMillis = 850
const SeekTurnSpeed = 4000
const SeekTurnMillis = 1500

const BackTurn1SpeedOuter = MaxSpeed
const BackTurn1SpeedInner = MaxSpeed / 2
const BackTurn1Millis = 400
const BackTurn2Speed = 5000
const BackTurn2Millis = 800
const BackMoveSpeed = MaxSpeed
const BackMoveMillis = 500
const BackTurn3Speed = 5000
const BackTurn3Millis = 1000

const CircleFindBorderMillis = 200
const CircleFindBorderOuterSpeed = MaxSpeed * 70 / 100
const CircleFindBorderInnerSpeed = MaxSpeed * 25 / 100
const CircleFindBorderOuterSpeedSlow = MaxSpeed * 31 / 100
const CircleFindBorderInnerSpeedSlow = MaxSpeed * 18 / 100

const CircleMillis = 2000
const CircleOuterSpeed = MaxSpeed
const CircleInnerSpeedLeft = 5300
const CircleInnerSpeedRight = 5300

const CircleAdjustInnerMax = 440
const CircleSpiralMillis = 500
const CircleSpiralOuterSpeed = MaxSpeed
const CircleSpiralInnerSpeed = 2000

const GoForwardMillis = 1000
const GoForwardSpeed = MaxSpeed

// const GoForwardTurnMillis = 600
const GoForwardTurnMillis = 0
const GoForwardTurnOuterSpeed = MaxSpeed
const GoForwardTurnInnerSpeed = 1000

const TurnBackPreMoveSpeed = MaxSpeed
const TurnBackPreMoveMillis = 200
const TurnBackMillis = 300
const TurnBackOuterSpeed = MaxSpeed
const TurnBackInnerSpeed = -MaxSpeed
const TurnBackMoveMillis = 650
const TurnBackMoveSpeed = MaxSpeed

const TrackFrontAngle = 15
const TrackSemiFrontAngle = 30
const TrackSemiFrontInnerSpeed = MaxSpeed / 2
const TrackMaxSpeed = MaxSpeed
const TrackOuterSpeed = 4000
const TrackInnerSpeed = -3000
const TrackSpeedReductionAngle = VisionMaxAngle - TrackFrontAngle
const TrackSpeedReductionMax = TrackOuterSpeed - TrackInnerSpeed
const TrackVisionIntensityIgnoreBorder = 40

// VisionSpeed is the speed of the eyes motor, max is 1560
const VisionSpeed = "450"

// const VisionFarValueFront = 90
// const VisionFarValueSide = 80

const VisionFarValueFront = 80
const VisionFarValueSide = 70
const VisionFarValueDelta = VisionFarValueFront - VisionFarValueSide

const VisionMaxAngle = (VisionMaxPosition * 9 / 25) + 45

const VisionMaxIntensity = 100

const VisionStartPosition = 150
const VisionStartPositionString = "150"
const VisionMaxPosition = 158
const VisionThresholdPosition = 155
const VisionEstimateReductionRange = 10
const VisionSpotWidth = 5 * 25 / 9
const VisionSpotSearchWidth = VisionMaxPosition - VisionSpotWidth

const VisionIgnoreBorderValue = 60
