#include<stdio.h>
#include<vector>
//#include <iostream>
#include<algorithm>
#include<string.h>
#include<stdio.h>

using namespace std;

extern "C" {
//extern void prints_l(char *str, unsigned int len);
int invoke()
{

	char buf[128];
	eosio::datastream<char *> ds(buf, sizeof(buf));

	int i = 0;
    	vector<int> v;

	prints_l("over0\n", 6);
	for( i = 0; i < 10; i++ )
	{
		v.push_back(i);
	}

	prints_l("over1\n", 6);
	//for( i = 0; i < v.size(); i++ )
	//{
	//	ds<< v[i] << " ";
	//}
	//printf("%s\n", ds);
	prints_l("over2\n", 6);

	printf("\n");
	//prints_l("over1\n", 6);


	printf("反序，从大到小\n");
	reverse(v.begin(), v.end());

	for (int i=0; i< v.size(); i++)
	{
		printf("%d ", v[i]);
	}

	//printf("\n");
	//prints_l("over2\n", 6);

	//printf("利用迭代器\n");
	vector<int>::iterator it;
	for (it = v.begin(); it != v.end(); it++)
	{
		//prints_l("over3\n", 6);
		printf("%d ", *it);
	}
	printf("\n");
	prints_l("over4\n", 6);

	return 0;
}
}
