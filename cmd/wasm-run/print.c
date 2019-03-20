#include<stdio.h>

int invoke()
{
	//__int128 a = 123;
	//__int128 b = 321;
	//__int128 c = a + b;
	unsigned long a = 1;
	unsigned long b = 2;
	unsigned long c = a + b;
	//long double a = 9.9999;
	//long double b = 8.9999;
	//long double c = a + b;

	//double b = 9.1234567;
	//unsigned int a = 11;
	//printf("hello world %d\n", a);
	fprintf(stderr, "hello world %d\n", c);
	//int b = 20;
	//char *ptr = "my name is steven";
	//printf("res = %x, res = %d, str = %s\n", a, b, ptr);
	return c;
}
