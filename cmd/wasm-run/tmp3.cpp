#include<stdlib.h>
#include<malloc.h>
#include<crypto.hpp>
#include<mydatastream.hpp>
#include<serialize.hpp>
#include<varint.hpp>

//typedef struct public_key {
//   unsigned_int        type;
//   std::array<char,33> data;
//   ONTLIB_SERIALIZE( public_key, (type)(data) )
//}public_key;
struct address {
	uint8_t bytes[20];
};

extern "C" {
void  ontio_assert( uint32_t test, const char* msg  )
{
	if ( test )
	{
		printf("%s\n", msg);
		abort();
	}
	
}


int invoke()
{
	//ontio::public_key t;
	//address tt;

	//size_t size = 128;
	//constexpr size_t max_stack_buffer_size = 512;
	//void *buffer = max_stack_buffer_size < size ? malloc(size) : alloca(size);
	//ontio::datastream<const char*> ds((char*)buffer, size);
	int a = 5;
	int b = 2;
	int c = a / b;
	printf("%d\n", c);
		


	//ds << t;
	/*
	ds << a;
	ds.seekp(0);
	ds >> b;
	if (a == b){
		printf("equal\n");
	} else {
		printf("not equal\n");
	}
	printf("0x%x\n", a);
	printf("0x%x\n", b);
	*/
	return 0;

}
}
