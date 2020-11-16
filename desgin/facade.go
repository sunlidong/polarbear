package desgin
import(
"fmt"

)

func NewAPI() API {
	return $apiImpl{
		a: NewAModuleAPI(),
		b: NewBModuleAPI(),
	}
}

type API interface{	
	Test() string
}

type apiImple struct{
	a AModuleAPI 
	b BModuleAPI
}

func (a *apiImple) Test() string {
	aRet :=a.a.TestA()
	bRet :=a.b.TestB()
	return fmt.Sprintf("%s\n%s",aRet,bRet)
}

func NewAModuleAPI() AModuleAPI {
	return &aModuleImpl{}
}

type AModuleAPI interface{
 	TestA() string
}

type aModuleImpl struct{
}

func(*aModuleImpl)TestA()string{
	return "a moudule running"
}

func NewBMoudleAPI() BModuleAPI{
	return &bModuleImpl{}
}
type BModuleAPI interface{
	TestB() string
}

type bModuleImpl struct{
}

func (*bModuleImpl) TestB()string{
	return "B moudule running"
}