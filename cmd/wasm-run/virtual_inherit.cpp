#include<iostream>
using namespace std;
class A
{
public:
    A(int v=100):X(v){};
    virtual void fooA(void){ cout<<"A::fooA"<<endl;}
    int X;
};
 
class B :virtual public A
{
public:
    B(int v=10):Y(v),A(100){};
    virtual void fooB(void){ cout<<"B::fooB"<<endl;}
    int Y;
};
 
class C : virtual public A
{
public:
    C(int v=20):Z(v),A(100){}
    virtual void fooC(void){ cout<<"C::fooC"<<endl;}
    int Z;
};

class XX:virtual public B{
public:
	XX(int v = 222):m_xx(v),B(10){}
	virtual void fooXX(void){}
	int m_xx;
};

class TMP_B:public A{
public:
	TMP_B(int v = 222):m_tmpb(v),A(100){}
	virtual void fooTMPB(void){cout<<"TMP_B::fooTMPB"<<endl;}
	int m_tmpb;
};

class TMP_D:public TMP_B,public C{
public:
	TMP_D(int v = 333):m_tmpd(v){}
	virtual void fooTMPD(void){cout<<"TMPD::fooTMPD()"<<endl;}
	int m_tmpd;
};

class TMP:virtual public A
{
public:
	TMP(int v = 111):tmp(v),A(100){}
	virtual void fooTMP(void){  cout<<"TMP:fooTMP"<<endl;}
	int tmp;
};

 
class D : public B, public C
{
public:
    D(int v =40):B(10),C(20),A(100),L(v){}
    virtual void fooD(void){ cout<<"D::fooD"<<endl; }
    int L;
};

class E : virtual public B, public C
{
public:
    E(int v =40):B(10),C(20),e(v){}
    virtual void fooE(void) { cout<<"E::fooE"<<endl;}
    int e;
};

class F : virtual public B, public C, public TMP
{
public:
    F(int v =40):B(10),C(20),A(100),f(v){}
    virtual void fooF(void){ cout<<"F::fooF"<<endl;}
    int f;
};

class G : virtual public B, public C, virtual public TMP
{
public:
    G(int v =40):B(10),C(20),A(100),g(v){}
    virtual void fooG(void){ cout<<"G::fooG"<<endl;}
    int g;
};

class X{
public:
	X(int v = 1000):m(v){}
	virtual void fooX(void){}
	int m;
	
};
class Y{
public:
	Y(int v = 1001):n(v){}
	virtual void fooY(void){}
	int n;
};


class Z:  public X, virtual public Y{
public:
	Z(int v = 1002):X(1000),Y(1001),o(v){}
	virtual void fooZ(void){}
	int o;
};


extern "C"{
int invoke()
{
	cout << "A" <<"   "<< sizeof(A) <<endl;
	cout << "B" <<"   "<< sizeof(B) <<endl;
	cout << "C" <<"   "<< sizeof(C) <<endl;
	cout << "D" <<"   "<< sizeof(D) <<endl;
	cout << "E" <<"   "<< sizeof(E) <<endl;
	cout << "Z" <<"   "<< sizeof(Z) <<endl;
	cout << "XX" <<"   "<< sizeof(XX) <<endl;
	return 0;
}
}


 /*
 
int main()
{
   
    A a;
    int *ptr;
	cout<<"start for class A---------------------"<<endl;
    ptr = (int*)&a;
    cout <<"start addr:"<< ptr << " sizeof = " << sizeof(a) <<endl<<endl;
	int i;
    for( i=0;i<sizeof(A)/sizeof(int);i++)
    {
		cout<<(ptr+i)<<'\t';
        if(ptr[i] < 10000)
        {
             cout << dec << ptr[i]<<endl;
        }
        else cout << hex << ptr[i] <<" = " << hex << * ((int*)(ptr[i])) <<endl;
    }
 




	cout<<endl<<endl<<endl;
    cout << "start for class B------------------------" <<endl;
 
    B b;
    ptr = (int*)&b;
    cout <<"start addr:" << ptr << " sizeof = " << sizeof(b) <<endl;
    for( i=0;i<sizeof(B)/sizeof(int);i++)
    {
		cout<<(ptr+i)<<'\t';
        if(ptr[i] < 10000)
        {
             cout << dec << ptr[i]<<endl;
        }
        else cout << hex << ptr[i] <<" = " << hex << * ((int*)(ptr[i])) <<endl;
    }
 




	cout<<endl<<endl<<endl;
    cout << "start for class D-------------------------" <<endl;
   
    D d;
    ptr = (int*)&d;
    cout <<"start addr:" << ptr << " sizeof = " << sizeof(d) <<endl<<endl;
    for( i=0;i<sizeof(D)/sizeof(int);i++)
    {
		cout<<(ptr+i)<<'\t';
        if(ptr[i] < 10000)
        {
             cout << dec << ptr[i]<<endl;
        }
        else cout << hex << ptr[i] <<" = " << hex << * ((int*)(ptr[i])) <<endl;
    }

    typedef void (*fun_ptr)(void);
	fun_ptr p_fun;
	int *VPTR_B_D = (int *)*ptr;          //注意这样调用，一般无法正确传递this指针，小心访问
	int *VPTR_C = (int *)*(ptr+3);
	int *VPTR_A = (int *)*(ptr+7);
	int *VPTR = VPTR_A;  


	p_fun = (fun_ptr)(*VPTR);
    p_fun();
    
	p_fun = (fun_ptr)(*(VPTR+1));
    p_fun();

	
	p_fun = (fun_ptr)(*(VPTR_D+2));
    p_fun();
	
	int *VBPTR_B = (int *)*(ptr+1);
	int *VBPTR_C = (int *)*(ptr+4);
	int *VBPTR = VBPTR_C;

	int j;
	for(j = 0;j < 2;j++)
	{
		cout<<hex<<"*(VBPTR+"<<j<<") = "<<*(VBPTR+j)<<endl;
	}
	cout<<dec;




	cout<<endl<<endl<<endl;
    cout << "start for class E-------------------------" <<endl;
   
    E e;
    ptr = (int*)&e;
    cout <<"start addr:" << ptr << " sizeof = " << sizeof(e) <<endl<<endl;
    for( i=0;i<sizeof(E)/sizeof(int);i++)
    {
		cout<<(ptr+i)<<'\t';
        if(ptr[i] < 10000)
        {
             cout << dec << ptr[i]<<endl;
        }
        else cout << hex << ptr[i] <<" = " << hex << * ((int*)(ptr[i])) <<endl;
    }
   


	int *VPTR_C_E = (int *)*ptr;          //注意这样调用，一般无法正确传递this指针，小心访问
	 VPTR_A = (int *)*(ptr+4);
    int *VPTR_B = (int *)*(ptr+6);
	 VPTR = VPTR_C_E; 
	 VPTR = VPTR_A;
	 VPTR = VPTR_B;
	p_fun = (fun_ptr)*(VPTR+0);
	p_fun();
	
    p_fun = (fun_ptr)*(VPTR+1);
	p_fun();
	


	int *VBPTR_C_E = (int *)*(ptr+1);
	VBPTR_B = (int *)*(ptr+7);
	VBPTR = VBPTR_C_E;
	VBPTR = VBPTR_B;

	for(j = 0;j < 3;j++)
	{
		cout<<hex<<"*(VBPTR+"<<j<<") = "<<*(VBPTR+j)<<endl;
	}
	cout<<dec;






	cout<<endl<<endl<<endl;
    cout << "start for class F-------------------------" <<endl;
   
    F f;
    ptr = (int*)&f;
    cout <<"start addr:" << ptr << " sizeof = " << sizeof(f) <<endl<<endl;
    for( i=0;i<sizeof(F)/sizeof(int);i++)
    {
		cout<<(ptr+i)<<'\t';
        if(ptr[i] < 10000)
        {
             cout << dec << ptr[i]<<endl;
        }
        else cout << hex << ptr[i] <<" = " << hex << * ((int*)(ptr[i])) <<endl;
    }


	int *VPTR_C_F = (int *)*ptr;          //注意这样调用，一般无法正确传递this指针，小心访问
	VPTR_A = (int *)*(ptr+7);
    VPTR_B = (int *)*(ptr+9);
	int *VPTR_TMP = (int *)*(ptr+3);
	VPTR = VPTR_C_F; 
	VPTR = VPTR_A;
	VPTR = VPTR_B;
	VPTR = VPTR_TMP;
	p_fun = (fun_ptr)*(VPTR+0);
	p_fun();

	int *VBPTR_C_F = (int *)*(ptr+1);
	VBPTR_B = (int *)*(ptr+10);
	int *VBPTR_TMP = (int *)*(ptr+4);
	VBPTR = VBPTR_C_F;
	VBPTR = VBPTR_B;
	VBPTR = VBPTR_TMP;

	for(j = 0;j < 8;j++)
	{
		cout<<hex<<"*(VBPTR+"<<j<<") = "<<*(VBPTR+j)<<endl;
	}
	cout<<dec;






	cout<<endl<<endl<<endl;
    cout << "start for class G-------------------------" <<endl;
   
    G g;
    ptr = (int*)&g;
    cout <<"start addr:" << ptr << " sizeof = " << sizeof(G) <<endl<<endl;
    for( i=0;i<sizeof(G)/sizeof(int);i++)
    {
		cout<<(ptr+i)<<'\t';
        if(ptr[i] < 10000)
        {
             cout << dec << ptr[i]<<endl;
        }
        else cout << hex << ptr[i] <<" = " << hex << * ((int*)(ptr[i])) <<endl;
    }
	

	int *VPTR_C_G = (int *)*ptr;          //注意这样调用，一般无法正确传递this指针，小心访问
	VPTR_A = (int *)*(ptr+4);
    VPTR_B = (int *)*(ptr+6);
    VPTR_TMP = (int *)*(ptr+9);
	VPTR = VPTR_C_G; 
	VPTR = VPTR_A;
	VPTR = VPTR_B;
	VPTR = VPTR_TMP;
	p_fun = (fun_ptr)*(VPTR);
	p_fun();


	int *VBPTR_C_G = (int *)*(ptr+1);
	VBPTR_B = (int *)*(ptr+7);
	VBPTR_TMP = (int *)*(ptr+10);
	VBPTR = VBPTR_C_G;
	VBPTR = VBPTR_B;
	VBPTR = VBPTR_TMP;

	for(j = 0;j < 8;j++)
	{
		cout<<hex<<"*(VBPTR+"<<j<<") = "<<*(VBPTR+j)<<endl;
	}
	cout<<dec;




	cout<<endl<<endl<<endl;
    cout << "start for class TMP_D-------------------------" <<endl;
   
    TMP_D tmp_d;
    ptr = (int*)&tmp_d;
    cout <<"start addr:" << ptr << " sizeof = " << sizeof(tmp_d) <<endl<<endl;
    for( i=0;i<sizeof(TMP_D)/sizeof(int);i++)
    {
		cout<<(ptr+i)<<'\t';
        if(ptr[i] < 10000)
        {
             cout << dec << ptr[i]<<endl;
        }
        else cout << hex << ptr[i] <<" = " << hex << * ((int*)(ptr[i])) <<endl;
    }

	cout<<"tmp_d.m_a = "<<tmp_d.X<<endl;



	int *VPTR_TMPBD = (int *)*ptr;          //注意这样调用，一般无法正确传递this指针，小心访问
	VPTR_A = VPTR_TMPBD;
	int *VPTR_A_virtual = (int *)*(ptr+7);
    int *VPTR_TMPB = VPTR_TMPBD;
    VPTR_C = (int *)*(ptr+3);



	VPTR = VPTR_TMPBD; 
	VPTR = VPTR_A_virtual;
//	VPTR = VPTR_C;
	p_fun = (fun_ptr)*(VPTR);
	p_fun();

	VBPTR_C = (int *)*(ptr+4);
	VBPTR = VBPTR_C;

	for(j = 0;j < 8;j++)
	{
		cout<<hex<<"*(VBPTR+"<<j<<") = "<<*(VBPTR+j)<<endl;
	}
	cout<<dec;



    return 0;
}
*/
