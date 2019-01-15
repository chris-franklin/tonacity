package tonacity

var scaleOrderNames = map[uint8]string{
	1: "Monotonic",
	2: "Ditonic", // Not to be confused with Diatonic scales, which are actually all heptatonic
	3: "Tritonic",
	4: "Tetratonic",
	5: "Pentatonic",
	6: "Hexatonic",
	7: "Heptatonic",
	8: "Octatonic",
	// 9,10,11 apparently not a thing
	12: "Chromatic", // aka dodecatonic
}
