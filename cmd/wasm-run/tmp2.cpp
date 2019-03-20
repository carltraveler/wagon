#include<stdlib.h>
#include<malloc.h>
#include<mydatastream.hpp>


/*
typedef unsigned long long   uint64_t;
template<typename T>
class mydatastream{
   public:
      mydatastream( T start, size_t s )
      :_start(start),_pos(start),_end(start+s){}
      inline void skip( size_t s ){ _pos += s; }
      inline bool read( char* d, size_t s ) {
        //ontio::check( size_t(_end - _pos) >= (size_t)s, "read" );
        memcpy( d, _pos, s );
        _pos += s;
        return true;
      }
      inline bool write( const char* d, size_t s ) {
        //ontio::check( _end - _pos >= (int32_t)s, "write" );
        memcpy( (void*)_pos, d, s );
        _pos += s;
        return true;
      }
      inline bool put(char c) {
        //ontio::check( _pos < _end, "put" );
        *_pos = c;
        ++_pos;
        return true;
      }

      inline bool get( unsigned char& c ) { return get( *(char*)&c ); }

      inline bool get( char& c )
      {
        //ontio::check( _pos < _end, "get" );
        c = *_pos;
        ++_pos;
        return true;
      }

      T pos()const { return _pos; }
      inline bool valid()const { return _pos <= _end && _pos >= _start;  }

      inline bool seekp(size_t p) { _pos = _start + p; return _pos <= _end; }

      inline size_t tellp()const      { return size_t(_pos - _start); }

      inline size_t remaining()const  { return _end - _pos; }
    private:
      T _start;
      T _pos;
      T _end;
};
*/


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
	uint64_t a = 0x20;
	uint64_t b = 1;
	size_t size = 128;
	constexpr size_t max_stack_buffer_size = 512;
	void *buffer = max_stack_buffer_size < size ? malloc(size) : alloca(size);
	ontio::datastream<const char*> ds((char*)buffer, size);
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
	return 0;

}
}
