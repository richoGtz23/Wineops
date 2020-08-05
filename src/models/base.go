package models

type Model interface {
	cypher() string
}

type BaseModel struct{
	id int
	Uuid string
	name string
}

func BaseModel()  {
	
}