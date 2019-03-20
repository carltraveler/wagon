#include<stdio.h>
#include<ontiolib/datastream.hpp>
#include<ontiolib/ontio.hpp>
#include<string>

using namespace ontio;

class hello : public contract {
public:
	using contract::contract;
	//[[eosio::action]]
	//[[ chenglin::xxx jdljfalksjflajl jdkjfakj ]]
	void hi(name user) {
		printf("hello world xxxx\n");
	}

	void hii(std::string s) {
		printf("hello world steven %s\n", s.c_str());
	}

	void xiiiiiqicheng(unsigned_int &a) {
		uint32_t b = 0x1234567;
		if (a == b)
			printf("hello world chenglin_hi 0x%x\n", a);
		else{
			printf("wrong args passed\n");
		}
	}

};

ONTIO_DISPATCH( hello, (hi)(hii)(xiiiiiqicheng) )

//extern "C"{
//void apply(char *str)
//{
//	printf("hello world. %s\n", str);
//}
//}
//


std::vector<char> HexToBytes(const std::string& hex) {
  std::vector<char> bytes;

  for (unsigned int i = 0; i < hex.length(); i += 2) {
    std::string byteString = hex.substr(i, 2);
    //printf("%s\n", byteString.c_str());
    uint8_t byte = (uint8_t) strtol(byteString.c_str(), NULL, 16);
    printf("%02x ", byte);
    bytes.push_back(byte);
  }
  printf("over\n");

  return bytes;
}

//std::array<uint8_t, ADDR_LEN>
template<size_t Size>
std::array<uint8_t,Size> HexToBytes(const std::string& hex) {
  std::array<uint8_t, Size> bytes;

  for (unsigned int i = 0; i < hex.length(); i += 2) {
    std::string byteString = hex.substr(i, 2);
    //printf("%s\n", byteString.c_str());
    uint8_t byte = (uint8_t) strtol(byteString.c_str(), NULL, 16);
    printf("%02x ", byte);
    bytes[i/2] = byte;
  }

  printf("over\n");

  return bytes;
}

extern "C"{


void  ontio_assert( uint32_t test, const char* msg  )
{
	if ( not test )
	{
		printf("%s\n", msg);
		abort();
	}
	
}

void print_byte(uint8_t *buffer, uint64_t len)
{
	uint64_t i = 0;
	for(; i < len; i++)
	{
		printf("%02X ",buffer[i]);
	}
}

void invoke()
{
	NativeArg tt;
	tt.version = 0x1;
	tt.method = std::string("transfer");
	//tt.from.bytes = HexToBytes(std::string("a36963b39d3eb14ddc1a9d016a14256ac594f8a4")).data();
	//tt.to.bytes = HexToBytes(std::string("a36963b39d3eb14ddc1a9d016a14256ac594f8a4")).data();
	
	//uint8_t *ttt = (uint8_t *)HexToBytes<20>(std::string("a36963b39d3eb14ddc1a9d016a14256ac594f8a4")).data();
	tt.from.bytes = HexToBytes<20>(std::string("a36963b39d3eb14ddc1a9d016a14256ac594f8a4"));
	tt.to.bytes = HexToBytes<20>(std::string("b0dc4eca16f06404201b3889b6adc2a6ed0246f4"));
	tt.value = 0x1234;

	//int i  = 0;
	//while(i < 20) {
	//	printf("%02x\n", tt.to.bytes.data()[i]);
	//	i++;
	//}
	//unsigned_int a = 0x12345678;
	//Address to;
	//to.bytes = HexToBytes<20>(std::string("b0dc4eca16f06404201b3889b6adc2a6ed0246f4"));
	
		
	//first save the runtime input buffer.
	//uint64_t action = name("hi").value;
	//uint64_t action = name("xiiiiiqicheng").value;
	//uint64_t action = name("hii").value;
	//unsigned_int a = 0x12345678;
	//uint32_t a = 0x78;
	//unsigned_int a = 0x12345678;
	//uint64_t a = 0x1234;
	//uint64_t a = 20;
	//name user("steven");
	//std::string s("now try string.");


	size_t size = pack_size(tt);
	//size_t size = 68;
	printf("size %lu\n", size);
	//printf("%lu\n",name( BOOST_PP_STRINGIZE("hi")  ).value);


	constexpr size_t max_stack_buffer_size = 512;
	void *buffer = max_stack_buffer_size < size ? malloc(size) : alloca(size);
	datastream<char*> ds((char*)buffer, size);
	
	//printf("%x\n", action);
	//ds << a;
	//WriteVarUint(ds, a);
	//WriteNativeArg<char *>(ds, tt);
	//WriteNativeArg(ds, tt);

	ds << tt;
	printf("ds bytes %lu\n", ds.tellp());
	print_byte((uint8_t *)buffer, ds.tellp());
	//ds << s;
	save_input_arg(buffer, size);

	ds.seekp(0);

	printf("ds bytes %lu\n", ds.remaining());
	NativeArg tt0;
	ds >> tt0;
	printf("version: %u\n", tt0.version);
	printf("method: %s\n", tt0.method.c_str());

	printf("from:  ");
	for(auto i : tt0.from.bytes){printf("%02x",i);}
	printf("\n");

	printf("to: ");
	for(auto i : tt0.to.bytes){printf("%02x",i);}
	printf("\n");

	printf("ontasset: 0x%x\n", tt0.value);

	throw "my thow";

	//unsigned_int b;
	//ds >> b;
	//printf("a value 0x%x\n", a.value);
	//printf("b value 0x%x\n", b.value);

	//apply(); //simulation entry call
	free(buffer);
}
}
