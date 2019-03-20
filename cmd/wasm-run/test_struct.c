#include<string.h>
#include<malloc.h>
//typedef unsigned int uint32_t;
//typedef int int32_t;
//typedef unsigned int size_t;
struct addr {
	unsigned int a;
	unsigned int b;
};

//struct addr myaddr(int a, int b, int c)
//{
//	struct addr t;
//	t.a = a + b + c;
//	t.b = b + b + c + 1;
//	return t;
//}

//unsigned int invoke(struct addr)
//{
//	return 
//}

int add(int a)
{
	return a + 1;
}

struct addr myaddr(int a)
{
	struct addr addr_t;
	memset(&addr_t, 0, sizeof(struct addr));
	addr_t.a = a + 99;
	addr_t.b = a + add(22);

	//struct addr *ptr_t = (struct addr *) malloc(sizeof(struct addr));
	//memset(ptr_t, 0, sizeof(struct addr));
	return addr_t;
}

//void *myaddr(int a)
//{
//	//struct addr addr_t;
//	int b = 9;
//	//memset(&addr_t, 0, sizeof(struct addr));
//	//addr_t.a = a + 99;
//	//addr_t.b = a + add(22);
//
//	//struct addr *ptr_t = (struct addr *) malloc(sizeof(struct addr));
//	//memset(ptr_t, 0, sizeof(struct addr));
//	return &b;
//}

//struct addr invoke(int a)
//struct addr invoke()
int invoke()
{
	struct addr addr_t = myaddr(66);
	//int a = (int)myaddr(66);
	return addr_t.b;
	//return a;
}
