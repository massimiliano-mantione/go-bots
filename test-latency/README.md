# Small program to test latency when reading EV3 sensors.

The program is used on the robot shown [here](https://www.youtube.com/watch?v=Ds_pBSqUnmU).

The rotating plates are seen by the IR sensor, which turns on the leds when it has a plate in front of it.
The color sensor reads the ambient light, and therefore the leds.

The program is simply an endless loop that reads the sensors and logs everything (the motor is in `run-direct` mode so nothing must be done about it).
The program also sleeps for one microsecond (1us) at the end of each iteration, to allow the scheduler to do its job.

Each loop iteration ideally takes less than one millisecond (1ms).
In normal conditions (without using `nice`) the program actually uses only about 20% of the available CPU time, and lots of iterations take more than 20 or 50 ms (because the program is preempted).

Decreasing the `nice` level iteration times get smaller, with more and more iterations taking the ideal 300us (or something similar).

The problem is that it takes *more* than one loop iteration for the information to propagate to the color sensor.

What I mean is: the color sensor reads the light from the leds a few iterations later than when the program turns on the leds.

Paradoxically, increasing the program priority (decreasing the `nice` level) makes things *worse*.

The files `data-n0X.txt` contain the iteration data logged at different niceness levels (`nice -n -X`).

The question is: is there some kernel thread or something similar that needs to run alongside the program so that the sensor data is properly available?
Would there be a way to "tune" the relative priorities of this thread and the program?

Note: the `go` garbage collector has nothing to do with this, it is never invoked in these runs.
