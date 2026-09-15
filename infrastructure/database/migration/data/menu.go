package data

import (
	"fp-kpl/infrastructure/database/schema"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func GetMenus(db *gorm.DB) []schema.Menu {
	var categories []schema.Category

	if err := db.Find(&categories).Error; err != nil {
		log.Fatalf("could not fetch categories: %v", err)
	}
	if len(categories) == 0 {
		log.Fatalf("no categories found")
	}

	categoryMap := make(map[string]string)
	for _, category := range categories {
		categoryMap[category.Name] = category.ID.String()
	}

	return []schema.Menu{
		{
			CategoryID:  getCategoryID(categoryMap, "Stand Indomie & Warkop"),
			Name:        "Indomie Goreng",
			ImageURL:    "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRSI6ixikF0ssei0fcZPF4LOEavsgcotLW1Ahy_HMzSPpPbLFdUsHBs65in&s=10",
			Price:       decimal.NewFromInt(6000),
			IsAvailable: true,
			CookingTime: schema.Duration{Duration: 5 * time.Minute},
			Description: "Indomie goreng dan sawi segar.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Stand Indomie & Warkop"),
			Name:        "Indomie Sedap Aceh Isi Dua",
			ImageURL:    "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQAXk7SyCNE8eNiudNlmGc2BehBRg3pCsOFCQBrWMZDwHh2n9hWc_b4SEPJ&s=10",
			Price:       decimal.NewFromInt(8000),
			IsAvailable: true,
			CookingTime: schema.Duration{Duration: 5 * time.Minute},
			Description: "Indomie kuah rasa soto hangat ditambah irisan cabai rawit segar.",
		},

		// 2. Stand Ayam Geprek
		{
			CategoryID:  getCategoryID(categoryMap, "Stand Ayam Geprek"),
			Name:        "Nasi Ayam Geprek Level Siswa",
			ImageURL:    "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQXo3QKT9nY2INbeoTUzLuwDWQ9y93-pj0-dUur4wjhqy8LfBf_wHxhXFw&s=10",
			Price:       decimal.NewFromInt(12000),
			IsAvailable: true,
			CookingTime: schema.Duration{Duration: 7 * time.Minute},
			Description: "Nasi hangat dengan ayam crispy yang digeprek sambal bawang .",
		},

		// 3. Stand Soto & Aneka Menu
		{
			CategoryID:  getCategoryID(categoryMap, "Stand Soto & Aneka Menu"),
			Name:        "Soto Ayam Kuah Kuning", // Menambahkan tanda petik penutup yang hilang
			ImageURL:    "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRRh9E_8sw6SMDLgn0b2adOaFcB7UQWQxPc1d_FnIzOIqrHpoXPU-oNfLw&s=10",
			Price:       decimal.NewFromInt(5000),
			IsAvailable: true,
			CookingTime: schema.Duration{Duration: 5 * time.Minute},
			Description: "Soto kuah kuning hangat kaya rempah dengan suwiran ayam, bihun, dan koya.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Stand Soto & Aneka Menu"),
			Name:        "usus crispy balado",
			ImageURL:    "https://laz-img-sg.alicdn.com/p/130fe67ce168efa50ac70103b017930d.jpg",
			Price:       decimal.NewFromInt(5000),
			IsAvailable: true,
			CookingTime: schema.Duration{Duration: 3 * time.Minute},
			Description: "ayam crispy dengan bumbu balado yang digoreng garing renyah.",
		},

		// 4. Stand Minuman & Cemilan
		{
			CategoryID:  getCategoryID(categoryMap, "Stand Minuman & Cemilan"),
			Name:        "Es tejus gula batu",
			ImageURL:    "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQT0XKOwlbgw8fQycfwpKNZhuA6lptocKwwuUvCjN1yS8LjgyQJ-EsbV8Y&s=10",
			Price:       decimal.NewFromInt(2000),
			IsAvailable: true,
			CookingTime: schema.Duration{Duration: 3 * time.Minute},
			Description: "Dibuat langsung di sekolah segar dan dingin.",
		},
		{
			CategoryID:  getCategoryID(categoryMap, "Stand Minuman & Cemilan"),
			Name:        "Wafer nabati keju",
			ImageURL:    "https://img.lazcdn.com/g/p/2b6cf21b7f86e42f123ff4a6aec1e11b.jpg_960x960q80.jpg_.webp",
			Price:       decimal.NewFromInt(2000),
			IsAvailable: true,
			CookingTime: schema.Duration{Duration: 3 * time.Minute},
			Description: "wafer nabati berbalu keju yg enak",
		},
	}
}

func getCategoryID(m map[string]string, name string) uuid.UUID {
	idStr, ok := m[name]
	if !ok {
		log.Fatalf("category %s not found", name)
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Fatalf("invalid category ID: %v", err)
	}

	return id
}