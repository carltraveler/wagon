extern "C" int invoke(void) {
	int a = 1;
	switch (a) {
		case 1:
			return a + 20;
		case 2:
			return a + 30;
		case 3:
			return a + 40;
		case 4:
			return a + 50;
		case 5:
			return a + 60;
		default:
			return a + 100;
	}
	return a;
}
