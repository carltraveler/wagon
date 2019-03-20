#include<stdio.h>
extern "C" {
int invoke()
{
	int a = 2;
	10 >> a;
	return a;
}
}
