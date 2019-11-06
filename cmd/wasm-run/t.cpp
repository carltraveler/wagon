extern "C" int func(int b) {
	int a = 8;
	return a + b;
}

typedef int (*F)(int);

extern "C" int invoke() {
	F f = func;
	return f(2);
}
