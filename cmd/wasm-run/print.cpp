#include<iostream>
#include<datastream.hpp>
//#include<stdio.h>
using namespace std;

extern "C" {
	int invoke()
	{
		printf("hello world\n");
		int a = 10;
		int b = 20;
		char *ptr = "my name is steven";
		printf("res = %x, res = %d, str = %s\n", a, b, ptr);
		eosio::datastream<<"hello world"<<endl;
		return 0;
	}
}

/*
int main()
{
	printf("hello world\n");
	int a = 10;
	int b = 20;
	char *ptr = "my name is steven";
	printf("res = %x, res = %d, str = %s\n", a, b, ptr);
	cout<<"hello world"<<endl;
}
*/
