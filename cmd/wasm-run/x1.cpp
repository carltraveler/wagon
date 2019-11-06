
extern "C" int invoke(void) {
	int a, b, c, d, e, f;
	a = 8 ; b = 9; c = 10; d = 11;
	while (true) {
		if (a == 8) {
			if (b == 9) {
				if (c == 10) {
					if (d == 11) {
						f = 10000;
						d = 12;
						continue;
					}
					f = f + 1;
					break;
				}
			}
		}
	}

	return f;
}

