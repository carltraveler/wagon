#include<stdio.h>
class MyOutstream
{
    public:
    const MyOutstream& operator << (int value)const;//对整型变量的重载
    const MyOutstream& operator << (char*str)const;//对字符串型的重载
};

const MyOutstream& MyOutstream::operator<<(int value)const
{
    printf("%d",value);
    return* this;//注意这个返回……
}

const MyOutstream& MyOutstream::operator<<(char* str)const
{
    printf("%s",str);
    return* this;//同样，这里也留意一下……
}

MyOutstream myout;//随时随地为我们服务的全局对象myout

extern "C" {
int invoke()
{
    int a=2003;
    char* myStr="Hello,World!";
    myout << myStr << "\n";
    return 0;
}
}
