package models

type User struct {
	UType    int      `json:"uType"`
	Email    string   `json:"email" gorm:"primaryKey"`
	Password []byte   `json:"-"`
	UserType Usertype `gorm:"foreignkey:UType"`
}

type Usertype struct {
	UserType string `json:"usertype"`
	ID       int    `json:"id" gorm:"primary_key"`
}

//var sellers []Seller
//db.Joins("JOIN shops s on s.id = sellers.shop_id").
//Joins("JOIN shop_types st on st.id = s.shop_type_id").
//Preload("Shop.ShopType").
//Where("st.name IN (?)", []string{"Store1", "Store2"}).
//Find(&sellers)
