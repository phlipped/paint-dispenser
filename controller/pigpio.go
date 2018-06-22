package controller

// Implements all the logic for dealing with gopigpio so that it can
// be hidden from the rest of the controller logic


/*

Calculations for stepping:

A) Thread pitch: 1.25mm
B) Steps per revolution: 200
C) Micro step level: 16
D) Microsteps per revolution: (B x C) 3200
E) Microsteps per millimeter: 

*/
