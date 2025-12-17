package ds

type CarbonateAcid struct {
	ID uint `gorm:"primaryKey"`
	// здесь создаем Unique key, указывая общий uniqueIndex
	AcidID      uint `gorm:"not null;uniqueIndex:idx_carbonate_acid"`
	CarbonateID uint `gorm:"not null;uniqueIndex:idx_carbonate_acid"`

	Mass   float32  `gorm:"default:null"`
	Result *float32 `gorm:"default:null"`

	Acid      Acid      `gorm:"foreignKey:AcidID"`
	Carbonate Carbonate `gorm:"foreignKey:CarbonateID"`
}
