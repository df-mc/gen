package gen

// normalCurve holds 128 values making up the shape of half a normal distribution curve.
var normalCurve = [128]float64{
	1, 0.9635997, 0.9362827, 0.9130436, 0.89228165, 0.87324303,
	0.8555006, 0.8387836, 0.8229072, 0.8077383, 0.793177,
	0.7791461, 0.7655842, 0.7524416, 0.73967725, 0.7272569,
	0.7151515, 0.7033361, 0.69178915, 0.68049186, 0.6694277,
	0.658582, 0.6479418, 0.63749546, 0.6272325, 0.6171434,
	0.6072195, 0.5974532, 0.58783704, 0.5783647, 0.56903,
	0.5598274, 0.5507518, 0.54179835, 0.5329627, 0.52424055,
	0.5156282, 0.50712204, 0.49871865, 0.49041483, 0.48220766,
	0.4740943, 0.46607214, 0.4581387, 0.45029163, 0.44252872,
	0.43484783, 0.427247, 0.41972435, 0.41227803, 0.40490642,
	0.39760786, 0.3903808, 0.3832238, 0.37613547, 0.36911446,
	0.3621595, 0.35526937, 0.34844297, 0.34167916, 0.33497685,
	0.3283351, 0.3217529, 0.3152294, 0.30876362, 0.30235484,
	0.29600215, 0.28970486, 0.2834622, 0.2772735, 0.27113807,
	0.2650553, 0.25902456, 0.2530453, 0.24711695, 0.241239,
	0.23541094, 0.22963232, 0.2239027, 0.21822165, 0.21258877,
	0.20700371, 0.20146611, 0.19597565, 0.19053204, 0.18513499,
	0.17978427, 0.17447963, 0.1692209, 0.16400786, 0.15884037,
	0.15371831, 0.14864157, 0.14361008, 0.13862377, 0.13368265,
	0.12878671, 0.12393598, 0.119130544, 0.11437051, 0.10965602,
	0.104987256, 0.10036444, 0.095787846, 0.0912578, 0.08677467,
	0.0823389, 0.077950984, 0.073611505, 0.06932112, 0.06508058,
	0.06089077, 0.056752663, 0.0526674, 0.048636295, 0.044660863,
	0.040742867, 0.03688439, 0.033087887, 0.029356318,
	0.025693292, 0.022103304, 0.018592102, 0.015167298,
	0.011839478, 0.008624485, 0.005548995, 0.0026696292,
}

var linearCurve = [128]float64{0.0078125, 0.015625, 0.0234375, 0.03125, 0.0390625, 0.046875, 0.0546875, 0.0625, 0.0703125, 0.078125, 0.0859375, 0.09375, 0.1015625, 0.109375, 0.1171875, 0.125, 0.1328125, 0.140625, 0.1484375, 0.15625, 0.1640625, 0.171875, 0.1796875, 0.1875, 0.1953125, 0.203125, 0.2109375, 0.21875, 0.2265625, 0.234375, 0.2421875, 0.25, 0.2578125, 0.265625, 0.2734375, 0.28125, 0.2890625, 0.296875, 0.3046875, 0.3125, 0.3203125, 0.328125, 0.3359375, 0.34375, 0.3515625, 0.359375, 0.3671875, 0.375, 0.3828125, 0.390625, 0.3984375, 0.40625, 0.4140625, 0.421875, 0.4296875, 0.4375, 0.4453125, 0.453125, 0.4609375, 0.46875, 0.4765625, 0.484375, 0.4921875, 0.5, 0.5078125, 0.515625, 0.5234375, 0.53125, 0.5390625, 0.546875, 0.5546875, 0.5625, 0.5703125, 0.578125, 0.5859375, 0.59375, 0.6015625, 0.609375, 0.6171875, 0.625, 0.6328125, 0.640625, 0.6484375, 0.65625, 0.6640625, 0.671875, 0.6796875, 0.6875, 0.6953125, 0.703125, 0.7109375, 0.71875, 0.7265625, 0.734375, 0.7421875, 0.75, 0.7578125, 0.765625, 0.7734375, 0.78125, 0.7890625, 0.796875, 0.8046875, 0.8125, 0.8203125, 0.828125, 0.8359375, 0.84375, 0.8515625, 0.859375, 0.8671875, 0.875, 0.8828125, 0.890625, 0.8984375, 0.90625, 0.9140625, 0.921875, 0.9296875, 0.9375, 0.9453125, 0.953125, 0.9609375, 0.96875, 0.9765625, 0.984375, 0.9921875, 1}
